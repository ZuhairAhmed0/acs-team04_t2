package station

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/database"
)

// ErrIdempotencyConflict indicates an idempotency key was reused with different parameters.
var ErrIdempotencyConflict = errors.New("idempotency key already used with different request parameters")

// GetStationWallet finds or initializes the station wallet.
func GetStationWallet(stationID uint64) (*Wallet, error) {
	var wallet Wallet
	err := database.DB.Where("owner_type = ? AND owner_id = ?", "station", stationID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = Wallet{
				OwnerType:        "station",
				OwnerID:          stationID,
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

// GetStationBranches lists all branches belonging to a station.
func GetStationBranches(stationID uint64) ([]Branch, error) {
	var branches []Branch
	err := database.DB.Where("station_id = ?", stationID).Order("is_main DESC, id ASC").Find(&branches).Error
	return branches, err
}

// AddStationBranch creates a new branch.
func AddStationBranch(branch *Branch) error {
	return database.DB.Create(branch).Error
}

// GetBranchByIDAndStation gets a branch guaranteeing it belongs to the given station.
func GetBranchByIDAndStation(branchID uint64, stationID uint64) (*Branch, error) {
	var b Branch
	err := database.DB.Where("id = ? AND station_id = ?", branchID, stationID).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// UpdateBranchFuelStatus updates the fuel status column.
func UpdateBranchFuelStatus(branchID uint64, fuelStatus string) error {
	return database.DB.Model(&Branch{}).Where("id = ?", branchID).Update("fuel_status", fuelStatus).Error
}

// UpdateBranchCongestion updates the congestion percentage column.
func UpdateBranchCongestion(branchID uint64, percentage int) error {
	return database.DB.Model(&Branch{}).Where("id = ?", branchID).Update("congestion_percentage", percentage).Error
}

// GetStationTransactions lists all transactions where station was source or target.
func GetStationTransactions(stationID uint64, limit int, offset int, txType, status, fromDate, toDate string) ([]Transaction, error) {
	var txs []Transaction
	query := database.DB.Model(&Transaction{}).Where(
		"((from_type = 'station' AND from_id = ?) OR (to_type = 'station' AND to_id = ?))",
		stationID, stationID,
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

func findIdempotentWithdrawal(tx *gorm.DB, stationID uint64, idempotencyKey, requestHash string) (*CashWithdrawal, error) {
	if idempotencyKey == "" {
		return nil, nil
	}
	var existing CashWithdrawal
	err := tx.Where("station_id = ? AND idempotency_key = ?", stationID, idempotencyKey).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if existing.RequestHash != nil && *existing.RequestHash != requestHash {
		return nil, ErrIdempotencyConflict
	}
	return &existing, nil
}

// CreateCashWithdrawalIdempotent locks the station wallet, checks idempotency and available
// unreserved balance, and records a pending cash withdrawal atomically.
func CreateCashWithdrawalIdempotent(stationID uint64, amount int64, idempotencyKey, requestHash string) (*CashWithdrawal, error) {
	var result CashWithdrawal
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var wallet Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("owner_type = ? AND owner_id = ?", "station", stationID).
			First(&wallet).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet = Wallet{OwnerType: "station", OwnerID: stationID}
				if createErr := tx.Create(&wallet).Error; createErr != nil {
					return createErr
				}
			} else {
				return err
			}
		}

		if existing, err := findIdempotentWithdrawal(tx, stationID, idempotencyKey, requestHash); err != nil {
			return err
		} else if existing != nil {
			result = *existing
			return nil
		}

		var openAmount int64
		if err := tx.Model(&CashWithdrawal{}).
			Where("station_id = ? AND status IN (?, ?)", stationID, "pending", "approved").
			Select("COALESCE(SUM(amount), 0)").
			Scan(&openAmount).Error; err != nil {
			return err
		}

		if wallet.AvailableBalance-openAmount < amount {
			return errors.New("insufficient available balance")
		}

		result = CashWithdrawal{
			StationID:      stationID,
			IdempotencyKey: &idempotencyKey,
			RequestHash:    &requestHash,
			Amount:         amount,
			Status:         "pending",
		}
		return tx.Create(&result).Error
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			if existing, lookupErr := findIdempotentWithdrawal(database.DB, stationID, idempotencyKey, requestHash); lookupErr != nil {
				return nil, lookupErr
			} else if existing != nil {
				return existing, nil
			}
		}
		return nil, err
	}
	return &result, nil
}

// CreateCashWithdrawal records a cash withdrawal request.
func CreateCashWithdrawal(cw *CashWithdrawal) error {
	return database.DB.Create(cw).Error
}

// GetStationCashWithdrawals lists requests for this station.
func GetStationCashWithdrawals(stationID uint64) ([]CashWithdrawal, error) {
	var list []CashWithdrawal
	err := database.DB.Where("station_id = ?", stationID).Order("created_at DESC, id DESC").Find(&list).Error
	return list, err
}
