package user

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"wa/app/database"
)

// GetUserByID retrieves a user record by ID.
func GetUserByID(userID uint64) (*User, error) {
	var u User
	if err := database.DB.First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetUserWallet retrieves the user's wallet.
func GetUserWallet(userID uint64) (*Wallet, error) {
	var wallet Wallet
	err := database.DB.Where("owner_type = ? AND owner_id = ?", "user", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Auto create if missing
			wallet = Wallet{
				OwnerType:        "user",
				OwnerID:          userID,
				TotalBalance:     0,
				AllocatedBalance: 0,
				AvailableBalance: 0,
			}
			if createErr := database.DB.Create(&wallet).Error; createErr != nil {
				return nil, createErr
			}
			return &wallet, nil
		}
		return nil, err
	}
	return &wallet, nil
}

// GetUserTransactions retrieves transactions involving the user.
func GetUserTransactions(userID uint64, limit int, offset int, txType, status, fromDate, toDate string) ([]Transaction, error) {
	var txs []Transaction
	query := database.DB.Model(&Transaction{}).Where(
		"((from_type = 'user' AND from_id = ?) OR (to_type = 'user' AND to_id = ?))",
		userID, userID,
	)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != "" {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate != "" {
		if len(toDate) == 10 {
			query = query.Where("created_at < (?::date + INTERVAL '1 day')", toDate)
		} else {
			query = query.Where("created_at <= ?", toDate)
		}
	}

	query = query.Order("created_at DESC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&txs).Error
	return txs, err
}

// GetActiveStationsWithBranches queries active stations and their branches, optionally filtering by city.
func GetActiveStationsWithBranches(city string) ([]StationBranchItem, error) {
	var items []StationBranchItem
	query := database.DB.Table("branches").
		Select("stations.id as station_id, stations.name as station_name, branches.id as branch_id, branches.city as branch_city, branches.address as branch_address, branches.fuel_status as fuel_status, branches.congestion_percentage as congestion_percentage").
		Joins("JOIN stations ON stations.id = branches.station_id").
		Where("stations.status = ?", "active")

	if trimmedCity := strings.TrimSpace(city); trimmedCity != "" {
		query = query.Where("LOWER(branches.city) = LOWER(?)", trimmedCity)
	}

	query = query.Order("stations.name ASC, branches.is_main DESC, branches.id ASC")
	err := query.Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
