package transfers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/database"
)

const recentDuplicatePurchaseWindow = 10 * time.Minute

// ErrIdempotencyConflict indicates an idempotency key was reused with different parameters.
var ErrIdempotencyConflict = errors.New("idempotency key already used with different request parameters")

// ExecutionOutcome holds the created transaction and any loyalty points earned.
type ExecutionOutcome struct {
	Transaction       *Transaction
	PointsEarned      int64
	LoyaltyPoints     *int64
	IdempotentReplay  bool
	DuplicatePurchase bool
	PreviousReference *string
}

func userActorKey(userID uint64) string {
	return "user:" + strconv.FormatUint(userID, 10)
}

func findIdempotentTransaction(tx *gorm.DB, actorKey, idempotencyKey, requestHash string) (*Transaction, error) {
	if idempotencyKey == "" {
		return nil, nil
	}
	var existing Transaction
	err := tx.Where("actor_key = ? AND idempotency_key = ?", actorKey, idempotencyKey).First(&existing).Error
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

func lockOrCreateWallet(tx *gorm.DB, ownerType string, ownerID uint64, createIfMissing bool) (*Wallet, error) {
	var wallet Wallet
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("owner_type = ? AND owner_id = ?", ownerType, ownerID).
		First(&wallet).Error
	if err == nil {
		return &wallet, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if !createIfMissing {
		return nil, gorm.ErrRecordNotFound
	}
	wallet = Wallet{OwnerType: ownerType, OwnerID: ownerID}
	if err := tx.Create(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func outcomeFromExisting(existing *Transaction, user *User) *ExecutionOutcome {
	var pts *int64
	if user != nil && user.Role == "user" {
		p := user.LoyaltyPoints
		pts = &p
	}
	return &ExecutionOutcome{
		Transaction:      existing,
		PointsEarned:     0,
		LoyaltyPoints:    pts,
		IdempotentReplay: true,
	}
}

func findRecentDuplicatePurchase(tx *gorm.DB, userID, branchID uint64, amount, pointsUsed int64) (*string, error) {
	var recent Transaction
	cutoff := time.Now().Add(-recentDuplicatePurchaseWindow)
	err := tx.Where(
		"type = ? AND status = ? AND from_type = ? AND from_id = ? AND branch_id = ? AND amount = ? AND loyalty_points_used = ? AND created_at >= ?",
		"fuel_purchase", "completed", "user", userID, branchID, amount, pointsUsed, cutoff,
	).Order("created_at DESC").First(&recent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	ref := recent.Reference
	return &ref, nil
}

func awardLoyaltyPoints(tx *gorm.DB, user *User) (int64, *int64, error) {
	if user == nil || user.Role != "user" {
		return 0, nil, nil
	}
	if err := tx.Model(&User{}).
		Where("id = ? AND role = ?", user.ID, "user").
		UpdateColumn("loyalty_points", gorm.Expr("loyalty_points + ?", database.LoyaltyActivityPoints)).Error; err != nil {
		return 0, nil, err
	}
	var updated User
	if err := tx.Select("loyalty_points").First(&updated, user.ID).Error; err != nil {
		return 0, nil, err
	}
	pts := updated.LoyaltyPoints
	return database.LoyaltyActivityPoints, &pts, nil
}

// ExecuteTopup adds funds to a user's wallet, records the transaction, and awards loyalty points.
func ExecuteTopup(userID uint64, amount int64, ref, idempotencyKey, requestHash string) (*ExecutionOutcome, error) {
	actorKey := userActorKey(userID)
	var outcome ExecutionOutcome
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("user not found")
		}

		wallet, err := lockOrCreateWallet(tx, "user", userID, true)
		if err != nil {
			return err
		}

		if existing, err := findIdempotentTransaction(tx, actorKey, idempotencyKey, requestHash); err != nil {
			return err
		} else if existing != nil {
			outcome = *outcomeFromExisting(existing, &user)
			return nil
		}

		wallet.AvailableBalance += amount
		wallet.TotalBalance += amount
		if err := tx.Save(wallet).Error; err != nil {
			return err
		}

		toType := "user"
		txRecord := Transaction{
			Reference:      ref,
			ActorKey:       &actorKey,
			IdempotencyKey: &idempotencyKey,
			RequestHash:    &requestHash,
			Type:           "topup",
			Amount:         amount,
			ToType:         &toType,
			ToID:           &userID,
			ToName:         user.Name,
			Status:         "completed",
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		earned, totalPts, err := awardLoyaltyPoints(tx, &user)
		if err != nil {
			return err
		}
		outcome = ExecutionOutcome{
			Transaction:   &txRecord,
			PointsEarned:  earned,
			LoyaltyPoints: totalPts,
		}
		return nil
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			if existing, lookupErr := findIdempotentTransaction(database.DB, actorKey, idempotencyKey, requestHash); lookupErr != nil {
				return nil, lookupErr
			} else if existing != nil {
				var user User
				_ = database.DB.First(&user, userID).Error
				return outcomeFromExisting(existing, &user), nil
			}
		}
		return nil, err
	}
	return &outcome, nil
}

// ExecuteAllocation transfers available funds from user to a station's shared balance.
func ExecuteAllocation(userID uint64, stationID uint64, amount int64, ref, idempotencyKey, requestHash string) (*Transaction, error) {
	actorKey := userActorKey(userID)
	var txRecord Transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("user not found")
		}

		var station Station
		if err := tx.First(&station, stationID).Error; err != nil {
			return errors.New("station not found")
		}
		if station.Status != "active" {
			return errors.New("station account is not active")
		}

		// Lock station wallet before user wallet ('station' < 'user') for deadlock-free ordering.
		stationWallet, err := lockOrCreateWallet(tx, "station", stationID, true)
		if err != nil {
			return err
		}

		userWallet, err := lockOrCreateWallet(tx, "user", userID, false)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if existing, idempErr := findIdempotentTransaction(tx, actorKey, idempotencyKey, requestHash); idempErr != nil {
					return idempErr
				} else if existing != nil {
					txRecord = *existing
					return nil
				}
				return errors.New("insufficient available balance")
			}
			return err
		}

		if existing, err := findIdempotentTransaction(tx, actorKey, idempotencyKey, requestHash); err != nil {
			return err
		} else if existing != nil {
			txRecord = *existing
			return nil
		}

		if userWallet.AvailableBalance < amount {
			return errors.New("insufficient available balance")
		}

		// Adjust balances:
		// User available decreases, allocated increases. Total balance remains total = available + allocated.
		userWallet.AvailableBalance -= amount
		userWallet.AllocatedBalance += amount
		if err := tx.Save(userWallet).Error; err != nil {
			return err
		}

		// Station available and total increase
		stationWallet.AvailableBalance += amount
		stationWallet.TotalBalance += amount
		if err := tx.Save(stationWallet).Error; err != nil {
			return err
		}

		fromType := "user"
		toType := "station"
		txRecord = Transaction{
			Reference:      ref,
			ActorKey:       &actorKey,
			IdempotencyKey: &idempotencyKey,
			RequestHash:    &requestHash,
			Type:           "allocation",
			Amount:         amount,
			FromType:       &fromType,
			FromID:         &userID,
			FromName:       user.Name,
			ToType:         &toType,
			ToID:           &stationID,
			ToName:         station.Name,
			Status:         "completed",
		}
		return tx.Create(&txRecord).Error
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			if existing, lookupErr := findIdempotentTransaction(database.DB, actorKey, idempotencyKey, requestHash); lookupErr != nil {
				return nil, lookupErr
			} else if existing != nil {
				return existing, nil
			}
		}
		return nil, err
	}
	return &txRecord, nil
}

// ExecuteFuelPurchase records a branch-level fuel purchase with optional loyalty point spending
// and awards loyalty points on purchases that do not redeem points.
func ExecuteFuelPurchase(userID uint64, branchID uint64, walletCharge, pointsUsed, discountAmount int64, paymentSource, ref, idempotencyKey, requestHash string) (*ExecutionOutcome, error) {
	actorKey := userActorKey(userID)
	var outcome ExecutionOutcome
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var branch Branch
		if err := tx.First(&branch, branchID).Error; err != nil {
			return errors.New("branch not found")
		}
		if branch.FuelStatus == "empty" {
			return errors.New("selected branch is out of fuel")
		}

		var station Station
		if err := tx.First(&station, branch.StationID).Error; err != nil {
			return errors.New("station not found")
		}
		if station.Status != "active" {
			return errors.New("station account is not active")
		}

		// Lock station wallet -> user wallet -> user row for deadlock-free ordering.
		stationWallet, err := lockOrCreateWallet(tx, "station", station.ID, true)
		if err != nil {
			return err
		}

		userWallet, err := lockOrCreateWallet(tx, "user", userID, walletCharge == 0)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var user User
				if userErr := tx.First(&user, userID).Error; userErr != nil {
					return errors.New("user not found")
				}
				if existing, idempErr := findIdempotentTransaction(tx, actorKey, idempotencyKey, requestHash); idempErr != nil {
					return idempErr
				} else if existing != nil {
					outcome = *outcomeFromExisting(existing, &user)
					return nil
				}
				if paymentSource == "allocated" {
					return errors.New("insufficient allocated balance")
				}
				return errors.New("insufficient available balance")
			}
			return err
		}

		var user User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return errors.New("user not found")
		}

		if existing, err := findIdempotentTransaction(tx, actorKey, idempotencyKey, requestHash); err != nil {
			return err
		} else if existing != nil {
			outcome = *outcomeFromExisting(existing, &user)
			return nil
		}

		if pointsUsed > 0 {
			if user.Role != "user" {
				return errors.New("loyalty points are only available for user accounts")
			}
			if user.LoyaltyPoints < pointsUsed {
				return errors.New("insufficient loyalty points")
			}
			user.LoyaltyPoints -= pointsUsed
			if err := tx.Model(&User{}).Where("id = ?", user.ID).UpdateColumn("loyalty_points", user.LoyaltyPoints).Error; err != nil {
				return err
			}
		}

		if walletCharge > 0 {
			if paymentSource == "allocated" {
				if userWallet.AllocatedBalance < walletCharge {
					return errors.New("insufficient allocated balance")
				}
				if stationWallet.AvailableBalance < walletCharge {
					return errors.New("insufficient station balance")
				}
				userWallet.AllocatedBalance -= walletCharge
				userWallet.TotalBalance -= walletCharge
				if err := tx.Save(userWallet).Error; err != nil {
					return err
				}
			} else {
				if userWallet.AvailableBalance < walletCharge {
					return errors.New("insufficient available balance")
				}
				userWallet.AvailableBalance -= walletCharge
				userWallet.TotalBalance -= walletCharge
				if err := tx.Save(userWallet).Error; err != nil {
					return err
				}
				stationWallet.AvailableBalance += walletCharge
				stationWallet.TotalBalance += walletCharge
				if err := tx.Save(stationWallet).Error; err != nil {
					return err
				}
			}
		}

		previousRef, err := findRecentDuplicatePurchase(tx, userID, branch.ID, walletCharge, pointsUsed)
		if err != nil {
			return err
		}

		fromType := "user"
		toType := "station"
		bID := branch.ID
		bName := fmt.Sprintf("%s - %s", branch.City, branch.Address)
		txRecord := Transaction{
			Reference:         ref,
			ActorKey:          &actorKey,
			IdempotencyKey:    &idempotencyKey,
			RequestHash:       &requestHash,
			Type:              "fuel_purchase",
			Amount:            walletCharge,
			FromType:          &fromType,
			FromID:            &userID,
			FromName:          user.Name,
			ToType:            &toType,
			ToID:              &station.ID,
			ToName:            station.Name,
			BranchID:          &bID,
			BranchName:        &bName,
			LoyaltyPointsUsed: pointsUsed,
			LoyaltyDiscount:   discountAmount,
			Status:            "completed",
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		var earned int64
		var totalPts *int64
		if pointsUsed == 0 {
			earned, totalPts, err = awardLoyaltyPoints(tx, &user)
			if err != nil {
				return err
			}
		} else if user.Role == "user" {
			p := user.LoyaltyPoints
			totalPts = &p
		}

		outcome = ExecutionOutcome{
			Transaction:       &txRecord,
			PointsEarned:      earned,
			LoyaltyPoints:     totalPts,
			DuplicatePurchase: previousRef != nil,
			PreviousReference: previousRef,
		}
		return nil
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			if existing, lookupErr := findIdempotentTransaction(database.DB, actorKey, idempotencyKey, requestHash); lookupErr != nil {
				return nil, lookupErr
			} else if existing != nil {
				var user User
				_ = database.DB.First(&user, userID).Error
				return outcomeFromExisting(existing, &user), nil
			}
		}
		return nil, err
	}
	return &outcome, nil
}

// ExecuteStationTransfer shifts balance between two stations.
func ExecuteStationTransfer(userID uint64, fromStationID uint64, toStationID uint64, amount int64, ref, idempotencyKey, requestHash string) (*Transaction, error) {
	actorKey := userActorKey(userID)
	var txRecord Transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var fromStation, toStation Station
		if err := tx.First(&fromStation, fromStationID).Error; err != nil {
			return errors.New("station not found")
		}
		if err := tx.First(&toStation, toStationID).Error; err != nil {
			return errors.New("station not found")
		}
		if fromStation.Status != "active" || toStation.Status != "active" {
			return errors.New("station account is not active")
		}

		// Lock station wallets in ascending station ID order, then user wallet ('station' < 'user').
		var fromWallet, toWallet *Wallet
		var err error
		if fromStationID < toStationID {
			fromWallet, err = lockOrCreateWallet(tx, "station", fromStationID, false)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			toWallet, err = lockOrCreateWallet(tx, "station", toStationID, true)
			if err != nil {
				return err
			}
		} else {
			toWallet, err = lockOrCreateWallet(tx, "station", toStationID, true)
			if err != nil {
				return err
			}
			fromWallet, err = lockOrCreateWallet(tx, "station", fromStationID, false)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		userWallet, userWalletErr := lockOrCreateWallet(tx, "user", userID, false)
		if userWalletErr != nil && !errors.Is(userWalletErr, gorm.ErrRecordNotFound) {
			return userWalletErr
		}

		if existing, err := findIdempotentTransaction(tx, actorKey, idempotencyKey, requestHash); err != nil {
			return err
		} else if existing != nil {
			txRecord = *existing
			return nil
		}

		if userWallet == nil || userWallet.AllocatedBalance < amount {
			return errors.New("insufficient allocated balance")
		}
		if fromWallet == nil || fromWallet.AvailableBalance < amount {
			return errors.New("insufficient available balance")
		}

		fromWallet.AvailableBalance -= amount
		fromWallet.TotalBalance -= amount
		if err := tx.Save(fromWallet).Error; err != nil {
			return err
		}

		toWallet.AvailableBalance += amount
		toWallet.TotalBalance += amount
		if err := tx.Save(toWallet).Error; err != nil {
			return err
		}

		fromType := "station"
		toType := "station"
		txRecord = Transaction{
			Reference:      ref,
			ActorKey:       &actorKey,
			IdempotencyKey: &idempotencyKey,
			RequestHash:    &requestHash,
			Type:           "station_transfer",
			Amount:         amount,
			FromType:       &fromType,
			FromID:         &fromStationID,
			FromName:       fromStation.Name,
			ToType:         &toType,
			ToID:           &toStationID,
			ToName:         toStation.Name,
			Status:         "completed",
		}
		return tx.Create(&txRecord).Error
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			if existing, lookupErr := findIdempotentTransaction(database.DB, actorKey, idempotencyKey, requestHash); lookupErr != nil {
				return nil, lookupErr
			} else if existing != nil {
				return existing, nil
			}
		}
		return nil, err
	}
	return &txRecord, nil
}

// GetUserTransfers lists top-up, allocation, fuel purchase, and transfer transactions for a user.
func GetUserTransfers(userID uint64) ([]Transaction, error) {
	var txs []Transaction
	err := database.DB.Where(
		"type IN ('allocation', 'topup', 'fuel_purchase', 'station_transfer') AND ((from_type = 'user' AND from_id = ?) OR (to_type = 'user' AND to_id = ?))",
		userID, userID,
	).Order("created_at DESC, id DESC").Find(&txs).Error
	return txs, err
}

// GetTransferByID returns one specific transfer by ID if caller has access.
func GetTransferByID(id uint64, userID uint64) (*Transaction, error) {
	var t Transaction
	err := database.DB.Where(
		"id = ? AND ((from_type = 'user' AND from_id = ?) OR (to_type = 'user' AND to_id = ?))",
		id, userID, userID,
	).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}
