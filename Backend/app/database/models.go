package database

import (
	"time"
)

const (
	LoyaltyWelcomePoints   int64 = 10
	LoyaltyActivityPoints  int64 = 5
	LoyaltyPointsPerGallon int64 = 10
	LoyaltyPointValueSDG   int64 = 100
)

// User represents an individual user or administrator account.
type User struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"type:varchar(150);not null" json:"name"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"type:varchar(255);not null" json:"-"`
	Phone         string    `gorm:"type:varchar(30);not null" json:"phone"`
	Role          string    `gorm:"type:varchar(20);not null;default:'user'" json:"role"`     // 'user' or 'admin'
	Status        string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // 'active' or 'suspended'
	LoyaltyPoints int64     `gorm:"not null;default:0" json:"loyalty_points"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Station represents a fuel station company account.
type Station struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string     `gorm:"type:varchar(180);not null" json:"name"`
	Location          string     `gorm:"type:varchar(255);not null" json:"location"`
	BankAccountNumber string     `gorm:"type:varchar(100);not null" json:"bank_account_number"`
	BankAccountName   string     `gorm:"type:varchar(180);not null" json:"bank_account_name"`
	BankName          string     `gorm:"type:varchar(100);not null" json:"bank_name"`
	Email             string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash      string     `gorm:"type:varchar(255);not null" json:"-"`
	Status            string     `gorm:"type:varchar(30);not null;default:'pending'" json:"status"` // 'pending', 'active', 'rejected', 'revoked'
	ReviewedBy        *uint64    `gorm:"index" json:"reviewed_by,omitempty"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
	RejectionReason   *string    `gorm:"type:text" json:"rejection_reason,omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Reviewer *User    `gorm:"foreignKey:ReviewedBy;references:ID" json:"reviewer,omitempty"`
	Branches []Branch `gorm:"foreignKey:StationID" json:"branches,omitempty"`
}

// Branch represents a physical station branch location.
type Branch struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	StationID            uint64    `gorm:"not null;index" json:"station_id"`
	Station              *Station  `gorm:"foreignKey:StationID;references:ID" json:"station,omitempty"`
	City                 string    `gorm:"type:varchar(120);not null;index" json:"city"`
	Address              string    `gorm:"type:varchar(255);not null" json:"address"`
	IsMain               bool      `gorm:"default:false;not null" json:"is_main"`
	FuelStatus           string    `gorm:"type:varchar(20);not null;default:'full'" json:"fuel_status"` // 'full', 'low', 'empty'
	CongestionPercentage int       `gorm:"not null;default:0" json:"congestion_percentage"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Wallet stores the balance for a user or station account.
type Wallet struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerType        string    `gorm:"type:varchar(20);not null;index:idx_owner" json:"owner_type"` // 'user' or 'station'
	OwnerID          uint64    `gorm:"not null;index:idx_owner" json:"owner_id"`
	TotalBalance     int64     `gorm:"not null;default:0" json:"total_balance"`
	AllocatedBalance int64     `gorm:"not null;default:0" json:"allocated_balance"`
	AvailableBalance int64     `gorm:"not null;default:0" json:"available_balance"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Transaction records all monetary flows in the system.
type Transaction struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Reference      string    `gorm:"type:varchar(80);uniqueIndex;not null" json:"reference"`
	ActorKey       *string   `gorm:"type:varchar(60);uniqueIndex:idx_tx_actor_idempotency" json:"-"`
	IdempotencyKey *string   `gorm:"type:varchar(120);uniqueIndex:idx_tx_actor_idempotency" json:"idempotency_key,omitempty"`
	RequestHash    *string   `gorm:"type:varchar(64)" json:"-"`
	Type           string    `gorm:"type:varchar(40);not null" json:"type"` // 'topup', 'allocation', 'fuel_purchase', 'station_transfer', 'withdrawal'
	Amount         int64     `gorm:"not null" json:"amount"`
	FromType       *string   `gorm:"type:varchar(20)" json:"from_type,omitempty"` // 'user' or 'station'
	FromID         *uint64   `gorm:"index" json:"from_id,omitempty"`
	FromName       string    `gorm:"type:varchar(180);default:''" json:"from_name,omitempty"`
	ToType         *string   `gorm:"type:varchar(20)" json:"to_type,omitempty"` // 'user' or 'station'
	ToID           *uint64   `gorm:"index" json:"to_id,omitempty"`
	ToName         string    `gorm:"type:varchar(180);default:''" json:"to_name,omitempty"`
	BranchID          *uint64   `gorm:"index" json:"branch_id,omitempty"`
	BranchName        *string   `gorm:"type:varchar(180)" json:"branch_name,omitempty"`
	LoyaltyPointsUsed int64     `gorm:"not null;default:0" json:"loyalty_points_used,omitempty"`
	LoyaltyDiscount   int64     `gorm:"not null;default:0" json:"loyalty_discount,omitempty"`
	Status            string    `gorm:"type:varchar(30);not null;default:'completed'" json:"status"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// CashWithdrawal tracks station requests to withdraw cash.
type CashWithdrawal struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StationID         uint64     `gorm:"not null;index;uniqueIndex:idx_cw_station_idempotency" json:"station_id"`
	Station           *Station   `gorm:"foreignKey:StationID;references:ID" json:"station,omitempty"`
	IdempotencyKey    *string    `gorm:"type:varchar(120);uniqueIndex:idx_cw_station_idempotency" json:"idempotency_key,omitempty"`
	RequestHash       *string    `gorm:"type:varchar(64)" json:"-"`
	Amount            int64      `gorm:"not null" json:"amount"`
	Status            string     `gorm:"type:varchar(30);not null;default:'pending'" json:"status"` // 'pending', 'approved', 'rejected', 'completed'
	ExternalReference *string    `gorm:"type:varchar(120)" json:"external_reference,omitempty"`
	ReviewedBy        *uint64    `gorm:"index" json:"reviewed_by,omitempty"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
	RejectionReason   *string    `gorm:"type:text" json:"rejection_reason,omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Reviewer *User `gorm:"foreignKey:ReviewedBy;references:ID" json:"reviewer,omitempty"`
}

// RefreshToken is a hashed long-lived session token tied to an account.
type RefreshToken struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	UserID    string    `gorm:"type:text;index;not null" json:"user_id"` // format: "user:<id>" or "station:<id>"
	TokenHash string    `gorm:"uniqueIndex;not null" json:"token_hash"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
