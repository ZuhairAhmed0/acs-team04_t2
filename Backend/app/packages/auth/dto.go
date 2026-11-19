package auth

import (
	"errors"
	"strings"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/database"
)

// ----*Models of auth package (aliased to canonical database models)*----
type User = database.User
type Station = database.Station
type Branch = database.Branch
type Wallet = database.Wallet
type RefreshToken = database.RefreshToken

// ----*DTOs of auth package*----

// RegisterUserRequest holds fields for user registration.
type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func (r *RegisterUserRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	r.Phone = strings.TrimSpace(r.Phone)
	for _, ch := range r.Phone {
		if ch < '0' || ch > '9' {
			return errors.New("phone number must contain digits only")
		}
	}
	if len(r.Phone) > 10 {
		return errors.New("phone number must be at most 10 digits")
	}
	return nil
}

// RegisterStationRequest holds fields for station registration.
type RegisterStationRequest struct {
	Name              string `json:"name" validate:"optional"`
	StationName       string `json:"station_name" validate:"optional"`
	Location          string `json:"location"`
	BankAccountNumber string `json:"bank_account_number"`
	BankAccountName   string `json:"bank_account_name"`
	BankName          string `json:"bank_name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
}

func (r *RegisterStationRequest) Validate() error {
	if r.Name == "" && r.StationName != "" {
		r.Name = r.StationName
	}
	r.Name = strings.TrimSpace(r.Name)
	r.Location = strings.TrimSpace(r.Location)
	r.BankAccountNumber = strings.TrimSpace(r.BankAccountNumber)
	r.BankAccountName = strings.TrimSpace(r.BankAccountName)
	r.BankName = strings.TrimSpace(r.BankName)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	return nil
}

// StationRegisterData is returned upon successful station registration.
type StationRegisterData struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// RegisterStationResponse is the JSON response for station registration.
type RegisterStationResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    StationRegisterData `json:"data"`
}

// LoginRequest handles login for user, station, or admin accounts.
type LoginRequest struct {
	Email    string `json:"email" validate:"optional"`
	UserID   string `json:"userId" validate:"optional"`
	Password string `json:"password"`
}

// AccountResponse is the public profile of an authenticated user or station.
type AccountResponse struct {
	ID            uint64 `json:"id"`
	Role          string `json:"role"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	Phone         string `json:"phone,omitempty"`
	LoyaltyPoints *int64 `json:"loyalty_points,omitempty"`
	Location      string `json:"location,omitempty"`
	BankName      string `json:"bank_name,omitempty"`
	AccountNo     string `json:"bank_account_number,omitempty"`
}

// LoginResponse is returned upon successful login.
type LoginResponse struct {
	Success      bool             `json:"success"`
	AccessToken  string           `json:"accessToken,omitempty"`
	RefreshToken string           `json:"refreshToken,omitempty"`
	ExpiresIn    int64            `json:"expiresIn"`
	User         *AccountResponse `json:"user,omitempty"`
}

// RefreshRequest is used to refresh or revoke tokens.
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"optional"`
}

// RefreshResponse returns a new access token.
type RefreshResponse struct {
	Success     bool             `json:"success"`
	AccessToken string           `json:"accessToken,omitempty"`
	ExpiresIn   int64            `json:"expiresIn"`
	User        *AccountResponse `json:"user,omitempty"`
}

// MeResponse returns current profile details.
type MeResponse struct {
	Success bool             `json:"success"`
	Data    *AccountResponse `json:"data"`
}

// LogoutResponse confirms token revocation.
type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
