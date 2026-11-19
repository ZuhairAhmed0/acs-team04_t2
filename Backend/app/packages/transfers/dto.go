package transfers

import (
	"time"
	"wa/app/database"
)

type User = database.User
type Station = database.Station
type Branch = database.Branch
type Wallet = database.Wallet
type Transaction = database.Transaction

// TopupRequest holds parameters for topping up a wallet.
type TopupRequest struct {
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key" validate:"optional"`
}

// AllocateRequest holds parameters for allocating user balance to a station.
type AllocateRequest struct {
	StationID      uint64 `json:"station_id"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key" validate:"optional"`
}

// FuelPurchaseRequest holds parameters for purchasing fuel at a station branch.
type FuelPurchaseRequest struct {
	BranchID          uint64 `json:"branch_id"`
	Amount            int64  `json:"amount" validate:"optional"`
	PaymentSource     string `json:"payment_source" validate:"optional"`
	LoyaltyPointsUsed int64  `json:"loyalty_points_used" validate:"optional"`
	IdempotencyKey    string `json:"idempotency_key" validate:"optional"`
}

// StationTransferRequest holds parameters for station-to-station transfer.
type StationTransferRequest struct {
	FromStationID  uint64 `json:"from_station_id"`
	ToStationID    uint64 `json:"to_station_id"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key" validate:"optional"`
}

// TransferResult is the standard response returned after an operation.
type TransferResult struct {
	Reference         string  `json:"reference"`
	Amount            int64   `json:"amount"`
	TotalAmount       int64   `json:"total_amount,omitempty"`
	FromName          string  `json:"from_name"`
	ToName            string  `json:"to_name"`
	BranchName        *string `json:"branch_name"`
	Status            string  `json:"status"`
	PointsEarned      int64   `json:"points_earned,omitempty"`
	LoyaltyPointsUsed int64   `json:"loyalty_points_used,omitempty"`
	LoyaltyDiscount   int64   `json:"loyalty_discount,omitempty"`
	LoyaltyPoints     *int64  `json:"loyalty_points,omitempty"`
	IdempotentReplay  bool    `json:"idempotent_replay,omitempty"`
	DuplicatePurchase bool    `json:"duplicate_purchase,omitempty"`
	PreviousReference *string `json:"previous_reference,omitempty"`
}

// TransferDetail represents the full transaction detail.
type TransferDetail struct {
	ID                uint64    `json:"id"`
	Reference         string    `json:"reference"`
	Type              string    `json:"type"`
	Amount            int64     `json:"amount"`
	FromType          *string   `json:"from_type,omitempty"`
	FromID            *uint64   `json:"from_id,omitempty"`
	FromName          string    `json:"from_name"`
	ToType            *string   `json:"to_type,omitempty"`
	ToID              *uint64   `json:"to_id,omitempty"`
	ToName            string    `json:"to_name"`
	BranchID          *uint64   `json:"branch_id,omitempty"`
	BranchName        *string   `json:"branch_name,omitempty"`
	LoyaltyPointsUsed int64     `json:"loyalty_points_used,omitempty"`
	LoyaltyDiscount   int64     `json:"loyalty_discount,omitempty"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}
