package admin

import (
	"time"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/database"
)

type User = database.User
type Station = database.Station
type Branch = database.Branch
type Wallet = database.Wallet
type Transaction = database.Transaction
type CashWithdrawal = database.CashWithdrawal

// DashboardData holds the administrative metrics.
type DashboardData struct {
	UsersCount               int64 `json:"users_count"`
	RegisteredStationsCount  int64 `json:"registered_stations_count"`
	ApprovedStationsCount    int64 `json:"approved_stations_count"`
	PendingStationsCount     int64 `json:"pending_stations_count"`
	PendingWithdrawalsCount  int64 `json:"pending_withdrawals_count"`
}

// UpdateUserStatusRequest toggles between 'active' and 'suspended'.
type UpdateUserStatusRequest struct {
	Status string `json:"status"`
}

// UpdateUserRoleRequest toggles between 'user' and 'admin'.
type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

// RejectStationRequest provides optional reason for station rejection.
type RejectStationRequest struct {
	RejectionReason string `json:"rejection_reason" validate:"optional"`
}

// RejectWithdrawalRequest provides optional reason for withdrawal rejection.
type RejectWithdrawalRequest struct {
	RejectionReason string `json:"rejection_reason" validate:"optional"`
}

// CompleteWithdrawalRequest provides an optional external bank transfer reference.
type CompleteWithdrawalRequest struct {
	ExternalReference string `json:"external_reference" validate:"optional"`
}

// UserListItem is returned when querying users.
type UserListItem struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// WithdrawalListItem provides enriched withdrawal view.
type WithdrawalListItem struct {
	ID                uint64     `json:"id"`
	StationID         uint64     `json:"station_id"`
	StationName       string     `json:"station_name"`
	BankName          string     `json:"bank_name"`
	BankAccountNumber string     `json:"bank_account_number"`
	BankAccountName   string     `json:"bank_account_name"`
	Amount            int64      `json:"amount"`
	Status            string     `json:"status"`
	ExternalReference *string    `json:"external_reference,omitempty"`
	ReviewedBy        *uint64    `json:"reviewed_by,omitempty"`
	ReviewedByName    string     `json:"reviewed_by_name,omitempty"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
	RejectionReason   *string    `json:"rejection_reason,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}
