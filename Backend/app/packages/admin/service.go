package admin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// GetDashboard returns metrics for the admin dashboard.
func GetDashboard() (*DashboardData, error) {
	users, regStations, approvedStations, pendingStations, pendingWithdrawals, err := GetAdminDashboardCounts()
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		UsersCount:              users,
		RegisteredStationsCount: regStations,
		ApprovedStationsCount:   approvedStations,
		PendingStationsCount:    pendingStations,
		PendingWithdrawalsCount: pendingWithdrawals,
	}, nil
}

// GetUsers lists all users formatted for admin review.
func GetUsers() ([]UserListItem, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	items := make([]UserListItem, len(users))
	for i, u := range users {
		items[i] = UserListItem{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Phone:     u.Phone,
			Role:      u.Role,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		}
	}
	return items, nil
}

// ChangeUserStatus validates and updates user status ('active' or 'suspended').
func ChangeUserStatus(id uint64, status string) error {
	status = strings.TrimSpace(strings.ToLower(status))
	if status != "active" && status != "suspended" {
		return errors.New("invalid status (must be 'active' or 'suspended')")
	}
	err := UpdateUserStatus(id, status)
	if err != nil {
		return err
	}
	return nil
}

// ChangeUserRole validates and updates user role ('user' or 'admin').
func ChangeUserRole(id uint64, role string) error {
	role = strings.TrimSpace(strings.ToLower(role))
	if role != "user" && role != "admin" {
		return errors.New("invalid user role (must be 'user' or 'admin')")
	}
	err := UpdateUserRole(id, role)
	if err != nil {
		return err
	}
	return nil
}

// ListPendingStations returns stations awaiting approval.
func ListPendingStations() ([]Station, error) {
	stations, err := GetPendingStations()
	if err != nil {
		return nil, err
	}
	return stations, nil
}

// ListStations returns all registered stations.
func ListStations() ([]Station, error) {
	stations, err := GetAllStations()
	if err != nil {
		return nil, err
	}
	return stations, nil
}

// ApproveStationRegistration marks a station as active.
func ApproveStationRegistration(id uint64, adminID uint64) error {
	err := ApproveStation(id, adminID)
	if err != nil {
		return err
	}
	return nil
}

// RejectStationRegistration rejects a station registration.
func RejectStationRegistration(id uint64, adminID uint64, reason string) error {
	err := RejectStation(id, adminID, reason)
	if err != nil {
		return err
	}
	return nil
}

// RevokeStationRegistration revokes station approval.
func RevokeStationRegistration(id uint64) error {
	err := RevokeStation(id)
	if err != nil {
		return err
	}
	return nil
}

// RemoveStation deletes a station.
func RemoveStation(id uint64) error {
	err := DeleteStation(id)
	if err != nil {
		return err
	}
	return nil
}

// ListWithdrawals returns all withdrawal requests.
func ListWithdrawals() ([]WithdrawalListItem, error) {
	list, err := GetAllCashWithdrawals()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// ApproveWithdrawal marks request as approved.
func ApproveWithdrawal(id uint64, adminID uint64) error {
	err := ApproveCashWithdrawal(id, adminID)
	if err != nil {
		return err
	}
	return nil
}

// RejectWithdrawal marks request as rejected.
func RejectWithdrawal(id uint64, adminID uint64, reason string) error {
	err := RejectCashWithdrawal(id, adminID, reason)
	if err != nil {
		return err
	}
	return nil
}

// CompleteWithdrawal settles the withdrawal and updates the station wallet balance.
func CompleteWithdrawal(id uint64, externalReference string) error {
	err := CompleteCashWithdrawal(id, externalReference)
	if err != nil {
		return err
	}
	return nil
}

// ParseAdminID parses uint64 id from JWT sub.
func ParseAdminID(sub string) (uint64, error) {
	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid admin id: %w", err)
	}
	return id, nil
}
