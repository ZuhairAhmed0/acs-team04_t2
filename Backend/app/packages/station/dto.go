package station

import (
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/database"
)

type Station = database.Station
type Branch = database.Branch
type Wallet = database.Wallet
type Transaction = database.Transaction
type CashWithdrawal = database.CashWithdrawal

// BranchDashboardItem holds branch info shown on station dashboard.
type BranchDashboardItem struct {
	ID                   uint64 `json:"id"`
	City                 string `json:"city"`
	Address              string `json:"address,omitempty"`
	FuelStatus           string `json:"fuel_status"`
	CongestionPercentage int    `json:"congestion_percentage"`
}

// StationDashboardData represents the dashboard response payload.
type StationDashboardData struct {
	Balance            int64                 `json:"balance"`
	BranchesCount      int                   `json:"branches_count"`
	Branches           []BranchDashboardItem `json:"branches"`
	RecentTransactions []Transaction         `json:"recent_transactions"`
}

// AddBranchRequest holds city and address to add a branch.
type AddBranchRequest struct {
	City    string `json:"city"`
	Address string `json:"address"`
}

// UpdateFuelStatusRequest updates a branch's fuel availability status.
type UpdateFuelStatusRequest struct {
	FuelStatus string `json:"fuel_status"`
}

// UpdateCongestionRequest updates a branch's congestion percentage (0-100).
type UpdateCongestionRequest struct {
	CongestionPercentage int `json:"congestion_percentage" validate:"optional"`
}

// CashWithdrawalRequest submits a withdrawal amount.
type CashWithdrawalRequest struct {
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key" validate:"optional"`
}

// CashWithdrawalResult returns id, amount, and status.
type CashWithdrawalResult struct {
	ID     uint64 `json:"id"`
	Amount int64  `json:"amount"`
	Status string `json:"status"`
}
