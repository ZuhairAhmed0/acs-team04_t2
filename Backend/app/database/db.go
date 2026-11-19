package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// IsUniqueViolation reports whether err is a PostgreSQL unique-constraint violation (SQLSTATE 23505).
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

var (
	DB *gorm.DB
)

// InitDB opens the database, runs automigrations on the database models and
// sets the package-level DB handle.
func InitDB(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database with dsn %q: %w", dsn, err)
	}
	DB = db

	tables := []any{
		&User{},
		&Station{},
		&Branch{},
		&Wallet{},
		&Transaction{},
		&CashWithdrawal{},
		&RefreshToken{},
	}

	if err := db.AutoMigrate(tables...); err != nil {
		return fmt.Errorf("failed to automigrate models: %w", err)
	}

	constraintSQL := `
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_wallets_non_negative') THEN
			ALTER TABLE wallets
			ADD CONSTRAINT chk_wallets_non_negative
			CHECK (available_balance >= 0 AND allocated_balance >= 0 AND total_balance >= 0);
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_loyalty_non_negative') THEN
			ALTER TABLE users
			ADD CONSTRAINT chk_users_loyalty_non_negative
			CHECK (loyalty_points >= 0);
		END IF;
	END $$;`
	if err := db.Exec(constraintSQL).Error; err != nil {
		return fmt.Errorf("failed to ensure wallet balance constraints: %w", err)
	}
	return nil
}
