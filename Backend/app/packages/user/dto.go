package user

import (
	"time"
	"wa/app/database"
)

type User = database.User
type Wallet = database.Wallet
type Transaction = database.Transaction
type Station = database.Station
type Branch = database.Branch

// WalletData holds balance numbers and user loyalty points summary.
type WalletData struct {
	TotalBalance      int64 `json:"total_balance"`
	AllocatedBalance  int64 `json:"allocated_balance"`
	AvailableBalance  int64 `json:"available_balance"`
	LoyaltyPoints     int64 `json:"loyalty_points"`
	FreeGasGallons    int64 `json:"free_gas_gallons"`
	LoyaltyPointValue int64 `json:"loyalty_point_value"`
}

// TransactionItem represents a user transaction item.
type TransactionItem struct {
	ID                uint64    `json:"id"`
	Reference         string    `json:"reference"`
	Type              string    `json:"type"`
	Amount            int64     `json:"amount"`
	FromName          string    `json:"from_name"`
	ToName            string    `json:"to_name"`
	BranchName        *string   `json:"branch_name"`
	LoyaltyPointsUsed int64     `json:"loyalty_points_used,omitempty"`
	LoyaltyDiscount   int64     `json:"loyalty_discount,omitempty"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

// UserDashboardData is the dashboard payload.
type UserDashboardData struct {
	Wallet             WalletData        `json:"wallet"`
	RecentTransactions []TransactionItem `json:"recent_transactions"`
}

// StationBranchItem is returned when listing stations/branches for users.
type StationBranchItem struct {
	StationID            uint64 `json:"station_id"`
	StationName          string `json:"station_name"`
	BranchID             uint64 `json:"branch_id"`
	BranchCity           string `json:"branch_city"`
	BranchAddress        string `json:"branch_address"`
	FuelStatus           string `json:"fuel_status"`
	CongestionPercentage int    `json:"congestion_percentage"`
}
