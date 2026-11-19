package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"wa/app/database"
	"wa/app/security"
)

func userAccountResponse(user *User) *AccountResponse {
	resp := &AccountResponse{
		ID:     user.ID,
		Role:   user.Role,
		Name:   user.Name,
		Email:  user.Email,
		Phone:  user.Phone,
		Status: user.Status,
	}
	if user.Role == "user" {
		pts := user.LoyaltyPoints
		resp.LoyaltyPoints = &pts
	}
	return resp
}

// RegisterUser registers a regular user and creates their wallet.
func RegisterUser(req RegisterUserRequest) (*AccountResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" || req.Phone == "" {
		return nil, errors.New("missing email or password")
	}

	exists, err := EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email is already registered")
	}

	pwHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:          req.Name,
		Email:         req.Email,
		PasswordHash:  pwHash,
		Phone:         req.Phone,
		Role:          "user",
		Status:        "active",
		LoyaltyPoints: database.LoyaltyWelcomePoints,
	}

	if _, err := CreateUserWithWallet(user); err != nil {
		return nil, err
	}

	return userAccountResponse(user), nil
}

// RegisterStation registers a station in pending status and creates its main branch and shared wallet.
func RegisterStation(req RegisterStationRequest) (*RegisterStationResponse, error) {
	if req.Name == "" {
		req.Name = req.StationName
	}
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Location == "" ||
		req.BankAccountNumber == "" || req.BankAccountName == "" || req.BankName == "" {
		return nil, errors.New("missing email or password")
	}

	exists, err := EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email is already registered")
	}

	pwHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	station := &Station{
		Name:              req.Name,
		Location:          req.Location,
		BankAccountNumber: req.BankAccountNumber,
		BankAccountName:   req.BankAccountName,
		BankName:          req.BankName,
		Email:             req.Email,
		PasswordHash:      pwHash,
		Status:            "pending",
	}

	branch := &Branch{
		City:       req.Location,
		Address:    req.Location,
		IsMain:     true,
		FuelStatus: "full",
	}

	if _, err := CreateStationWithBranchAndWallet(station, branch); err != nil {
		return nil, err
	}

	return &RegisterStationResponse{
		Success: true,
		Message: "Station registration submitted",
		Data: StationRegisterData{
			ID:     station.ID,
			Name:   station.Name,
			Status: station.Status,
		},
	}, nil
}

// Login authenticates a user, admin, or station account.
func Login(req LoginRequest) (*LoginResponse, error) {
	identifier := strings.TrimSpace(req.Email)
	if identifier == "" {
		identifier = strings.TrimSpace(req.UserID)
	}
	if identifier == "" || strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("missing email or password")
	}

	// 1. Check users table (handles role "user" and "admin")
	user, err := GetUserByEmail(identifier)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if !security.CheckPasswordHash(req.Password, user.PasswordHash) {
			return nil, errors.New("invalid email or password")
		}
		if user.Status == "suspended" {
			return nil, errors.New("account is suspended or disabled")
		}

		sub := strconv.FormatUint(user.ID, 10)
		accessToken, err := security.GenerateToken(sub, user.Role, []string{user.Role})
		if err != nil {
			return nil, errors.New("error generating token")
		}

		plainRefresh, refreshHash, err := security.GenerateRefreshToken()
		if err != nil {
			return nil, errors.New("error generating token")
		}
		if err := CreateRefreshToken("user:"+sub, refreshHash, time.Now().Add(security.RefreshDuration)); err != nil {
			return nil, errors.New("error persisting token")
		}

		return &LoginResponse{
			Success:      true,
			AccessToken:  accessToken,
			RefreshToken: plainRefresh,
			ExpiresIn:    int64(security.AccessDuration.Seconds()),
			User:         userAccountResponse(user),
		}, nil
	}

	// 2. Check stations table (handles role "station")
	station, err := GetStationByEmail(identifier)
	if err != nil {
		return nil, err
	}
	if station != nil {
		if !security.CheckPasswordHash(req.Password, station.PasswordHash) {
			return nil, errors.New("invalid email or password")
		}
		if station.Status == "pending" {
			return nil, errors.New("station registration is pending admin approval")
		}
		if station.Status == "rejected" {
			return nil, errors.New("station registration was rejected")
		}
		if station.Status == "revoked" {
			return nil, errors.New("station registration has been revoked")
		}

		sub := strconv.FormatUint(station.ID, 10)
		accessToken, err := security.GenerateToken(sub, "station", []string{"station"})
		if err != nil {
			return nil, errors.New("error generating token")
		}

		plainRefresh, refreshHash, err := security.GenerateRefreshToken()
		if err != nil {
			return nil, errors.New("error generating token")
		}
		if err := CreateRefreshToken("station:"+sub, refreshHash, time.Now().Add(security.RefreshDuration)); err != nil {
			return nil, errors.New("error persisting token")
		}

		return &LoginResponse{
			Success:      true,
			AccessToken:  accessToken,
			RefreshToken: plainRefresh,
			ExpiresIn:    int64(security.AccessDuration.Seconds()),
			User: &AccountResponse{
				ID:        station.ID,
				Role:      "station",
				Name:      station.Name,
				Email:     station.Email,
				Status:    station.Status,
				Location:  station.Location,
				BankName:  station.BankName,
				AccountNo: station.BankAccountNumber,
			},
		}, nil
	}

	return nil, errors.New("invalid email or password")
}

// GetMe returns the profile of the authenticated caller.
func GetMe(userIDStr string, role string) (*MeResponse, error) {
	id, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token: %w", err)
	}

	if role == "station" {
		station, err := GetStationByID(id)
		if err != nil {
			return nil, err
		}
		if station == nil {
			return nil, errors.New("station not found")
		}
		return &MeResponse{
			Success: true,
			Data: &AccountResponse{
				ID:        station.ID,
				Role:      "station",
				Name:      station.Name,
				Email:     station.Email,
				Status:    station.Status,
				Location:  station.Location,
				BankName:  station.BankName,
				AccountNo: station.BankAccountNumber,
			},
		}, nil
	}

	user, err := GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return &MeResponse{
		Success: true,
		Data:    userAccountResponse(user),
	}, nil
}

// Logout invalidates a refresh token.
func Logout(plainRefreshToken string) error {
	plainRefreshToken = strings.TrimSpace(plainRefreshToken)
	if plainRefreshToken == "" {
		return errors.New("refresh token is required")
	}
	hash := security.HashToken(plainRefreshToken)
	return RevokeRefreshToken(hash)
}

// RefreshAccessToken issues a fresh access token from an active refresh token.
func RefreshAccessToken(plainRefreshToken string) (*RefreshResponse, error) {
	plainRefreshToken = strings.TrimSpace(plainRefreshToken)
	if plainRefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}
	hash := security.HashToken(plainRefreshToken)
	rt, err := GetValidRefreshToken(hash)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	parts := strings.SplitN(rt.UserID, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid or expired refresh token")
	}
	kind := parts[0]
	idStr := parts[1]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if kind == "station" {
		station, err := GetStationByID(id)
		if err != nil || station == nil || station.Status != "active" {
			return nil, errors.New("invalid or expired refresh token")
		}
		accessToken, err := security.GenerateToken(idStr, "station", []string{"station"})
		if err != nil {
			return nil, errors.New("error generating token")
		}
		return &RefreshResponse{
			Success:     true,
			AccessToken: accessToken,
			ExpiresIn:   int64(security.AccessDuration.Seconds()),
			User: &AccountResponse{
				ID:        station.ID,
				Role:      "station",
				Name:      station.Name,
				Email:     station.Email,
				Status:    station.Status,
				Location:  station.Location,
				BankName:  station.BankName,
				AccountNo: station.BankAccountNumber,
			},
		}, nil
	}

	user, err := GetUserByID(id)
	if err != nil || user == nil || user.Status != "active" {
		return nil, errors.New("invalid or expired refresh token")
	}
	accessToken, err := security.GenerateToken(idStr, user.Role, []string{user.Role})
	if err != nil {
		return nil, errors.New("error generating token")
	}
	return &RefreshResponse{
		Success:     true,
		AccessToken: accessToken,
		ExpiresIn:   int64(security.AccessDuration.Seconds()),
		User:        userAccountResponse(user),
	}, nil
}
