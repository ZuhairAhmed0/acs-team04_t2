package admin

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wa/app/database"
)

func generateReference(prefix string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	randomNum := n.Int64() + 100000
	return fmt.Sprintf("%s-%d-%06d", prefix, time.Now().Unix(), randomNum)
}

// GetAdminDashboardCounts aggregates system stats.
func GetAdminDashboardCounts() (users, regStations, approvedStations, pendingStations, pendingWithdrawals int64, err error) {
	if err = database.DB.Model(&User{}).Count(&users).Error; err != nil {
		return
	}
	if err = database.DB.Model(&Station{}).Count(&regStations).Error; err != nil {
		return
	}
	if err = database.DB.Model(&Station{}).Where("status = ?", "active").Count(&approvedStations).Error; err != nil {
		return
	}
	if err = database.DB.Model(&Station{}).Where("status = ?", "pending").Count(&pendingStations).Error; err != nil {
		return
	}
	err = database.DB.Model(&CashWithdrawal{}).Where("status = ?", "pending").Count(&pendingWithdrawals).Error
	return
}

// GetAllUsers returns all users ordered by creation date.
func GetAllUsers() ([]User, error) {
	var users []User
	err := database.DB.Order("created_at DESC").Find(&users).Error
	return users, err
}

// UpdateUserStatus activates or suspends a user.
func UpdateUserStatus(id uint64, status string) error {
	res := database.DB.Model(&User{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// UpdateUserRole changes user role between 'user' and 'admin'.
func UpdateUserRole(id uint64, role string) error {
	res := database.DB.Model(&User{}).Where("id = ?", id).Update("role", role)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetPendingStations returns stations awaiting approval.
func GetPendingStations() ([]Station, error) {
	var stations []Station
	err := database.DB.Preload("Branches").Where("status = ?", "pending").Order("created_at ASC").Find(&stations).Error
	return stations, err
}

// GetAllStations returns all stations with branches and reviewer info.
func GetAllStations() ([]Station, error) {
	var stations []Station
	err := database.DB.Preload("Branches").Preload("Reviewer").Order("created_at DESC").Find(&stations).Error
	return stations, err
}

// ApproveStation sets station status to active.
func ApproveStation(id uint64, adminID uint64) error {
	now := time.Now()
	res := database.DB.Model(&Station{}).Where("id = ?", id).Updates(map[string]any{
		"status":      "active",
		"reviewed_by": adminID,
		"reviewed_at": now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("station not found")
	}
	return nil
}

// RejectStation marks a station registration as rejected.
func RejectStation(id uint64, adminID uint64, reason string) error {
	now := time.Now()
	res := database.DB.Model(&Station{}).Where("id = ?", id).Updates(map[string]any{
		"status":           "rejected",
		"reviewed_by":      adminID,
		"reviewed_at":      now,
		"rejection_reason": reason,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("station not found")
	}
	return nil
}

// RevokeStation revokes an approved station's active status.
func RevokeStation(id uint64) error {
	res := database.DB.Model(&Station{}).Where("id = ?", id).Update("status", "revoked")
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("station not found")
	}
	return nil
}

// DeleteStation deletes a station and its branches, withdrawals, and wallet.
func DeleteStation(id uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("station_id = ?", id).Delete(&Branch{}).Error; err != nil {
			return err
		}
		if err := tx.Where("station_id = ?", id).Delete(&CashWithdrawal{}).Error; err != nil {
			return err
		}
		if err := tx.Where("owner_type = ? AND owner_id = ?", "station", id).Delete(&Wallet{}).Error; err != nil {
			return err
		}
		res := tx.Delete(&Station{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("station not found")
		}
		return nil
	})
}

// GetAllCashWithdrawals returns all withdrawal requests enriched with station and reviewer names.
func GetAllCashWithdrawals() ([]WithdrawalListItem, error) {
	var list []WithdrawalListItem
	err := database.DB.Table("cash_withdrawals").
		Select("cash_withdrawals.id, cash_withdrawals.station_id, stations.name as station_name, stations.bank_name, stations.bank_account_number, stations.bank_account_name, cash_withdrawals.amount, cash_withdrawals.status, cash_withdrawals.external_reference, cash_withdrawals.reviewed_by, users.name as reviewed_by_name, cash_withdrawals.reviewed_at, cash_withdrawals.rejection_reason, cash_withdrawals.created_at").
		Joins("JOIN stations ON stations.id = cash_withdrawals.station_id").
		Joins("LEFT JOIN users ON users.id = cash_withdrawals.reviewed_by").
		Order("cash_withdrawals.created_at DESC, cash_withdrawals.id DESC").
		Scan(&list).Error
	return list, err
}

// ApproveCashWithdrawal marks a withdrawal request as approved.
func ApproveCashWithdrawal(id uint64, adminID uint64) error {
	now := time.Now()
	res := database.DB.Model(&CashWithdrawal{}).Where("id = ? AND status = ?", id, "pending").Updates(map[string]any{
		"status":      "approved",
		"reviewed_by": adminID,
		"reviewed_at": now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("cash withdrawal request not found")
	}
	return nil
}

// RejectCashWithdrawal marks a withdrawal request as rejected.
func RejectCashWithdrawal(id uint64, adminID uint64, reason string) error {
	now := time.Now()
	res := database.DB.Model(&CashWithdrawal{}).Where("id = ? AND status != ?", id, "completed").Updates(map[string]any{
		"status":           "rejected",
		"reviewed_by":      adminID,
		"reviewed_at":      now,
		"rejection_reason": reason,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("cash withdrawal request not found")
	}
	return nil
}

// CompleteCashWithdrawal completes the withdrawal, deducts balance from the station wallet, and records the transaction.
func CompleteCashWithdrawal(id uint64, externalReference string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var cw CashWithdrawal
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&cw, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("cash withdrawal request not found")
			}
			return err
		}
		if cw.Status == "completed" {
			return errors.New("cash withdrawal request is already settled")
		}
		if cw.Status != "approved" {
			return errors.New("cash withdrawal request must be approved before completion")
		}

		var station Station
		if err := tx.First(&station, cw.StationID).Error; err != nil {
			return errors.New("station not found")
		}

		var wallet Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("owner_type = ? AND owner_id = ?", "station", cw.StationID).
			First(&wallet).Error; err != nil {
			return errors.New("insufficient available balance")
		}
		if wallet.AvailableBalance < cw.Amount {
			return errors.New("insufficient available balance")
		}

		wallet.AvailableBalance -= cw.Amount
		wallet.TotalBalance -= cw.Amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		cw.Status = "completed"
		if trimmedRef := strings.TrimSpace(externalReference); trimmedRef != "" {
			cw.ExternalReference = &trimmedRef
		}
		if err := tx.Save(&cw).Error; err != nil {
			return err
		}

		fromType := "station"
		transaction := Transaction{
			Reference: generateReference("WTH"),
			Type:      "withdrawal",
			Amount:    cw.Amount,
			FromType:  &fromType,
			FromID:    &cw.StationID,
			FromName:  station.Name,
			ToName:    station.BankName,
			Status:    "completed",
		}
		return tx.Create(&transaction).Error
	})
}
