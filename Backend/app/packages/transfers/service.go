package transfers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"wa/app/database"
	"wa/app/helpers"
)

func buildTransferResult(outcome *ExecutionOutcome) *TransferResult {
	tx := outcome.Transaction
	var totalAmount int64
	if tx.LoyaltyDiscount > 0 {
		totalAmount = tx.Amount + tx.LoyaltyDiscount
	}
	return &TransferResult{
		Reference:         tx.Reference,
		Amount:            tx.Amount,
		TotalAmount:       totalAmount,
		FromName:          tx.FromName,
		ToName:            tx.ToName,
		BranchName:        tx.BranchName,
		Status:            tx.Status,
		PointsEarned:      outcome.PointsEarned,
		LoyaltyPointsUsed: tx.LoyaltyPointsUsed,
		LoyaltyDiscount:   tx.LoyaltyDiscount,
		LoyaltyPoints:     outcome.LoyaltyPoints,
		IdempotentReplay:  outcome.IdempotentReplay,
		DuplicatePurchase: outcome.DuplicatePurchase,
		PreviousReference: outcome.PreviousReference,
	}
}

func generateReference(prefix string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	randomNum := n.Int64() + 100000
	return fmt.Sprintf("%s-%d-%06d", prefix, time.Now().Unix(), randomNum)
}

// Topup adds funds to the authenticated user's wallet.
func Topup(userID uint64, req TopupRequest) (*TransferResult, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	ref := generateReference("TOP")
	reqHash := helpers.HashIdempotencyPayload(fmt.Sprintf("topup:%d:%d", userID, req.Amount))
	outcome, err := ExecuteTopup(userID, req.Amount, ref, req.IdempotencyKey, reqHash)
	if err != nil {
		return nil, err
	}

	return buildTransferResult(outcome), nil
}

// Allocate moves balance from a user to a station.
func Allocate(userID uint64, req AllocateRequest) (*TransferResult, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if req.StationID == 0 {
		return nil, errors.New("station not found")
	}

	ref := generateReference("TRF")
	reqHash := helpers.HashIdempotencyPayload(fmt.Sprintf("allocate:%d:%d:%d", userID, req.StationID, req.Amount))
	tx, err := ExecuteAllocation(userID, req.StationID, req.Amount, ref, req.IdempotencyKey, reqHash)
	if err != nil {
		return nil, err
	}

	return &TransferResult{
		Reference:  tx.Reference,
		Amount:     tx.Amount,
		FromName:   tx.FromName,
		ToName:     tx.ToName,
		BranchName: tx.BranchName,
		Status:     tx.Status,
	}, nil
}

// PurchaseFuel processes a gas purchase at a specific station branch, optionally spending loyalty points
// (including 0 SDG points-only purchases), and awards loyalty points on wallet-paid purchases.
func PurchaseFuel(userID uint64, req FuelPurchaseRequest) (*TransferResult, error) {
	if req.Amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}
	if req.BranchID == 0 {
		return nil, errors.New("branch not found")
	}
	if req.LoyaltyPointsUsed < 0 {
		return nil, errors.New("loyalty_points_used cannot be negative")
	}
	source := strings.TrimSpace(strings.ToLower(req.PaymentSource))
	if source == "" {
		source = "available"
	}
	if source != "available" && source != "allocated" && source != "loyalty" {
		return nil, errors.New("invalid payment source")
	}

	var pointsUsed, discountAmount, walletCharge int64
	gallonValueSDG := database.LoyaltyPointsPerGallon * database.LoyaltyPointValueSDG
	if source == "loyalty" {
		if req.LoyaltyPointsUsed > 0 {
			if req.LoyaltyPointsUsed%database.LoyaltyPointsPerGallon != 0 {
				return nil, fmt.Errorf("loyalty_points_used must be a multiple of %d", database.LoyaltyPointsPerGallon)
			}
			expectedValue := req.LoyaltyPointsUsed * database.LoyaltyPointValueSDG
			if req.Amount > 0 && req.Amount != expectedValue {
				return nil, errors.New("loyalty_points_used does not match purchase amount")
			}
			pointsUsed = req.LoyaltyPointsUsed
			discountAmount = expectedValue
		} else {
			if req.Amount <= 0 {
				return nil, errors.New("loyalty_points_used or amount must be greater than zero")
			}
			if req.Amount%gallonValueSDG != 0 {
				return nil, fmt.Errorf("amount must be a multiple of %d SDG (%d loyalty points per gallon) when paying with loyalty points", gallonValueSDG, database.LoyaltyPointsPerGallon)
			}
			pointsUsed = req.Amount / database.LoyaltyPointValueSDG
			discountAmount = req.Amount
		}
		walletCharge = 0
	} else {
		if req.Amount == 0 && req.LoyaltyPointsUsed == 0 {
			return nil, errors.New("amount must be greater than zero")
		}
		if req.LoyaltyPointsUsed > 0 && req.LoyaltyPointsUsed%database.LoyaltyPointsPerGallon != 0 {
			return nil, fmt.Errorf("loyalty_points_used must be a multiple of %d", database.LoyaltyPointsPerGallon)
		}
		pointsUsed = req.LoyaltyPointsUsed
		discountAmount = pointsUsed * database.LoyaltyPointValueSDG
		if req.Amount == 0 {
			walletCharge = 0
		} else {
			if discountAmount > req.Amount {
				return nil, errors.New("loyalty discount cannot exceed purchase amount")
			}
			walletCharge = req.Amount - discountAmount
		}
	}

	ref := generateReference("PUR")
	reqHash := helpers.HashIdempotencyPayload(fmt.Sprintf("purchase:%d:%d:%d:%s:%d", userID, req.BranchID, req.Amount, source, pointsUsed))
	outcome, err := ExecuteFuelPurchase(userID, req.BranchID, walletCharge, pointsUsed, discountAmount, source, ref, req.IdempotencyKey, reqHash)
	if err != nil {
		return nil, err
	}

	return buildTransferResult(outcome), nil
}

// StationToStation transfers allocated balance between two stations.
func StationToStation(userID uint64, req StationTransferRequest) (*TransferResult, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if req.FromStationID == 0 || req.ToStationID == 0 {
		return nil, errors.New("station not found")
	}
	if req.FromStationID == req.ToStationID {
		return nil, errors.New("source and destination stations must be different")
	}

	ref := generateReference("STF")
	reqHash := helpers.HashIdempotencyPayload(fmt.Sprintf("station_transfer:%d:%d:%d:%d", userID, req.FromStationID, req.ToStationID, req.Amount))
	tx, err := ExecuteStationTransfer(userID, req.FromStationID, req.ToStationID, req.Amount, ref, req.IdempotencyKey, reqHash)
	if err != nil {
		return nil, err
	}

	return &TransferResult{
		Reference:  tx.Reference,
		Amount:     tx.Amount,
		FromName:   tx.FromName,
		ToName:     tx.ToName,
		BranchName: tx.BranchName,
		Status:     tx.Status,
	}, nil
}

func mapTransferDetail(t Transaction) TransferDetail {
	return TransferDetail{
		ID:                t.ID,
		Reference:         t.Reference,
		Type:              t.Type,
		Amount:            t.Amount,
		FromType:          t.FromType,
		FromID:            t.FromID,
		FromName:          t.FromName,
		ToType:            t.ToType,
		ToID:              t.ToID,
		ToName:            t.ToName,
		BranchID:          t.BranchID,
		BranchName:        t.BranchName,
		LoyaltyPointsUsed: t.LoyaltyPointsUsed,
		LoyaltyDiscount:   t.LoyaltyDiscount,
		Status:            t.Status,
		CreatedAt:         t.CreatedAt,
	}
}

// GetTransfers returns the list of transfer transactions for a user.
func GetTransfers(userID uint64) ([]TransferDetail, error) {
	txs, err := GetUserTransfers(userID)
	if err != nil {
		return nil, err
	}

	details := make([]TransferDetail, len(txs))
	for i, t := range txs {
		details[i] = mapTransferDetail(t)
	}
	return details, nil
}

// GetTransfer returns one transfer detail if accessible by the user.
func GetTransfer(id uint64, userID uint64) (*TransferDetail, error) {
	t, err := GetTransferByID(id, userID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.New("transfer not found")
	}

	detail := mapTransferDetail(*t)
	return &detail, nil
}

// ParseUserID parses user uint64 id from JWT sub.
func ParseUserID(sub string) (uint64, error) {
	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user id: %w", err)
	}
	return id, nil
}
