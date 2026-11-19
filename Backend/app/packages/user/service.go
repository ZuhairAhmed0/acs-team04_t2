package user

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"wa/app/database"
)

func buildWalletData(userID uint64, wallet *Wallet) (WalletData, error) {
	data := WalletData{
		TotalBalance:     wallet.TotalBalance,
		AllocatedBalance: wallet.AllocatedBalance,
		AvailableBalance: wallet.AvailableBalance,
	}
	u, err := GetUserByID(userID)
	if err != nil {
		return WalletData{}, err
	}
	if u != nil && u.Role == "user" {
		data.LoyaltyPoints = u.LoyaltyPoints
		data.FreeGasGallons = u.LoyaltyPoints / database.LoyaltyPointsPerGallon
		data.LoyaltyPointValue = database.LoyaltyPointValueSDG
	}
	return data, nil
}

func mapTransactions(txs []Transaction) []TransactionItem {
	items := make([]TransactionItem, len(txs))
	for i, t := range txs {
		items[i] = TransactionItem{
			ID:                t.ID,
			Reference:         t.Reference,
			Type:              t.Type,
			Amount:            t.Amount,
			FromName:          t.FromName,
			ToName:            t.ToName,
			BranchName:        t.BranchName,
			LoyaltyPointsUsed: t.LoyaltyPointsUsed,
			LoyaltyDiscount:   t.LoyaltyDiscount,
			Status:            t.Status,
			CreatedAt:         t.CreatedAt,
		}
	}
	return items
}

// GetDashboard returns the user dashboard overview.
func GetDashboard(userID uint64) (*UserDashboardData, error) {
	wallet, err := GetUserWallet(userID)
	if err != nil {
		return nil, err
	}

	walletData, err := buildWalletData(userID, wallet)
	if err != nil {
		return nil, err
	}

	txs, err := GetUserTransactions(userID, 3, 0, "", "", "", "")
	if err != nil {
		return nil, err
	}

	return &UserDashboardData{
		Wallet:             walletData,
		RecentTransactions: mapTransactions(txs),
	}, nil
}

// GetWallet returns the user's wallet balances.
func GetWallet(userID uint64) (*WalletData, error) {
	wallet, err := GetUserWallet(userID)
	if err != nil {
		return nil, err
	}
	data, err := buildWalletData(userID, wallet)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetStations returns the list of approved stations and their branches.
func GetStations(city string) ([]StationBranchItem, error) {
	items, err := GetActiveStationsWithBranches(city)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []StationBranchItem{}
	}
	return items, nil
}

// GetTransactions returns filtered transaction history for a user.
func GetTransactions(userID uint64, txType, status, fromDate, toDate string) ([]TransactionItem, error) {
	txs, err := GetUserTransactions(userID, 100, 0, txType, status, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	return mapTransactions(txs), nil
}

// ExportTransactionsCSV builds the CSV string for user transactions.
func ExportTransactionsCSV(userID uint64, txType, status, fromDate, toDate string) ([]byte, error) {
	txs, err := GetUserTransactions(userID, 0, 0, txType, status, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header matching README section 8
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

// ParseUserID parses uint64 id from JWT sub claim.
func ParseUserID(sub string) (uint64, error) {
	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user id: %w", err)
	}
	return id, nil
}
