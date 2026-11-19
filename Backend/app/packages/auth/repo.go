package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"wa/app/database"
)

// EmailExists checks if the given email already exists in either users or stations.
func EmailExists(email string) (bool, error) {
	var count int64
	if err := database.DB.Model(&User{}).Where("LOWER(email) = LOWER(?)", email).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	if err := database.DB.Model(&Station{}).Where("LOWER(email) = LOWER(?)", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateUserWithWallet creates a user and their initial wallet atomically.
func CreateUserWithWallet(user *User) (*Wallet, error) {
	var wallet Wallet
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		wallet = Wallet{
			OwnerType:        "user",
			OwnerID:          user.ID,
			TotalBalance:     0,
			AllocatedBalance: 0,
			AvailableBalance: 0,
		}
		return tx.Create(&wallet).Error
	})
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// CreateStationWithBranchAndWallet registers a station, creates its main branch, and its shared wallet atomically.
func CreateStationWithBranchAndWallet(station *Station, branch *Branch) (*Wallet, error) {
	var wallet Wallet
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(station).Error; err != nil {
			return err
		}
		branch.StationID = station.ID
		branch.IsMain = true
		if branch.FuelStatus == "" {
			branch.FuelStatus = "full"
		}
		if err := tx.Create(branch).Error; err != nil {
			return err
		}
		wallet = Wallet{
			OwnerType:        "station",
			OwnerID:          station.ID,
			TotalBalance:     0,
			AllocatedBalance: 0,
			AvailableBalance: 0,
		}
		return tx.Create(&wallet).Error
	})
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetUserByEmail finds a user by email.
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := database.DB.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetStationByEmail finds a station by email.
func GetStationByEmail(email string) (*Station, error) {
	var station Station
	err := database.DB.Where("LOWER(email) = LOWER(?)", email).First(&station).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &station, nil
}

// GetUserByID finds a user by numeric ID.
func GetUserByID(id uint64) (*User, error) {
	var user User
	err := database.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetStationByID finds a station by numeric ID.
func GetStationByID(id uint64) (*Station, error) {
	var station Station
	err := database.DB.First(&station, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &station, nil
}

// CreateRefreshToken records a hashed refresh token in the database.
func CreateRefreshToken(userID string, tokenHash string, expiresAt time.Time) error {
	record := &RefreshToken{
		ID:        tokenHash,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	return database.DB.Create(record).Error
}

// GetValidRefreshToken retrieves an active refresh token.
func GetValidRefreshToken(tokenHash string) (*RefreshToken, error) {
	var rt RefreshToken
	err := database.DB.Where("token_hash = ? AND revoked = ? AND expires_at > ?", tokenHash, false, time.Now()).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found or expired")
		}
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken marks a refresh token as revoked.
func RevokeRefreshToken(tokenHash string) error {
	res := database.DB.Model(&RefreshToken{}).Where("token_hash = ?", tokenHash).Update("revoked", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("refresh token revoked")
	}
	return nil
}
