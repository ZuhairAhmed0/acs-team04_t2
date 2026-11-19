package station

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"wa/app/helpers"
)

// GetDashboard aggregates data for the station dashboard page.
func GetDashboard(stationID uint64) (*StationDashboardData, error) {
	wallet, err := GetStationWallet(stationID)
	if err != nil {
		return nil, err
	}

	branches, err := GetStationBranches(stationID)
	if err != nil {
		return nil, err
	}

	txs, err := GetStationTransactions(stationID, 5, 0, "", "", "", "")
	if err != nil {
		return nil, err
	}

	branchItems := make([]BranchDashboardItem, len(branches))
	for i, b := range branches {
		branchItems[i] = BranchDashboardItem{
			ID:                   b.ID,
			City:                 b.City,
			Address:              b.Address,
			FuelStatus:           b.FuelStatus,
			CongestionPercentage: b.CongestionPercentage,
		}
	}

	return &StationDashboardData{
		Balance:            wallet.AvailableBalance,
		BranchesCount:      len(branches),
		Branches:           branchItems,
		RecentTransactions: txs,
	}, nil
}

// GetWallet returns the station wallet balances.
func GetWallet(stationID uint64) (*Wallet, error) {
	return GetStationWallet(stationID)
}

// GetTransactions lists transactions for this station.
func GetTransactions(stationID uint64, txType, status, fromDate, toDate string) ([]Transaction, error) {
	return GetStationTransactions(stationID, 100, 0, txType, status, fromDate, toDate)
}

// ExportTransactionsCSV builds the CSV export for station transactions.
func ExportTransactionsCSV(stationID uint64, txType, status, fromDate, toDate string) ([]byte, error) {
	txs, err := GetStationTransactions(stationID, 0, 0, txType, status, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"reference", "type", "amount", "from_name", "to_name", "branch_name", "status", "created_at"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, t := range txs {
		bName := ""
		if t.BranchName != nil {
			bName = *t.BranchName
		}
		row := []string{
			t.Reference,
			t.Type,
			strconv.FormatInt(t.Amount, 10),
			t.FromName,
			t.ToName,
			bName,
			t.Status,
			t.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetBranches returns all branches for this station.
func GetBranches(stationID uint64) ([]Branch, error) {
	return GetStationBranches(stationID)
}

// AddBranch creates a new branch for the station.
func AddBranch(stationID uint64, req AddBranchRequest) (*Branch, error) {
	city := strings.TrimSpace(req.City)
	address := strings.TrimSpace(req.Address)
	if city == "" || address == "" {
		return nil, errors.New("city and address are required")
	}

	branch := &Branch{
		StationID:  stationID,
		City:       city,
		Address:    address,
		IsMain:     false,
		FuelStatus: "full",
	}

	if err := AddStationBranch(branch); err != nil {
		return nil, err
	}
	return branch, nil
}

// UpdateFuelStatus updates the fuel status for a station's branch.
func UpdateFuelStatus(stationID uint64, branchID uint64, status string) error {
	status = strings.TrimSpace(strings.ToLower(status))
	if status != "full" && status != "low" && status != "empty" {
		return errors.New("invalid fuel status (must be 'full', 'low', or 'empty')")
	}

	branch, err := GetBranchByIDAndStation(branchID, stationID)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	return UpdateBranchFuelStatus(branchID, status)
}

// UpdateCongestion updates the congestion percentage (0-100) for a station's branch.
func UpdateCongestion(stationID uint64, branchID uint64, percentage int) error {
	if percentage < 0 || percentage > 100 {
		return errors.New("congestion_percentage must be between 0 and 100")
	}

	branch, err := GetBranchByIDAndStation(branchID, stationID)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	return UpdateBranchCongestion(branchID, percentage)
}

// SubmitCashWithdrawal validates balance and submits a pending withdrawal idempotently.
func SubmitCashWithdrawal(stationID uint64, req CashWithdrawalRequest) (*CashWithdrawalResult, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	reqHash := helpers.HashIdempotencyPayload(fmt.Sprintf("withdrawal:%d:%d", stationID, req.Amount))
	cw, err := CreateCashWithdrawalIdempotent(stationID, req.Amount, req.IdempotencyKey, reqHash)
	if err != nil {
		return nil, err
	}

	return &CashWithdrawalResult{
		ID:     cw.ID,
		Amount: cw.Amount,
		Status: cw.Status,
	}, nil
}

// GetCashWithdrawals lists requests submitted by the station.
func GetCashWithdrawals(stationID uint64) ([]CashWithdrawal, error) {
	return GetStationCashWithdrawals(stationID)
}

// ParseStationID parses station uint64 id from JWT sub.
func ParseStationID(sub string) (uint64, error) {
	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid station id: %w", err)
	}
	return id, nil
}
