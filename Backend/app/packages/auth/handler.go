package auth

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/helpers"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/security"
)

// HandleRegisterUser handles POST /api/v1/auth/register/user.
func HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	req, err := helpers.DecodeValid[RegisterUserRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := RegisterUser(req)
	if err != nil {
		switch err.Error() {
		case "email is already registered":
			helpers.RespondWithError(w, http.StatusConflict, err.Error())
		case "missing email or password":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	resp := map[string]any{
		"success": true,
		"message": "User registered successfully",
		"data":    user,
	}
	_ = helpers.Encode(w, r, http.StatusCreated, resp)
}

// HandleRegisterStation handles POST /api/v1/auth/register/station.
func HandleRegisterStation(w http.ResponseWriter, r *http.Request) {
	req, err := helpers.DecodeValid[RegisterStationRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := RegisterStation(req)
	if err != nil {
		switch err.Error() {
		case "email is already registered":
			helpers.RespondWithError(w, http.StatusConflict, err.Error())
		case "missing email or password":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusCreated, resp)
}

// HandleLogin handles POST /api/v1/auth/login.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	req, err := helpers.DecodeValid[LoginRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := Login(req)
	if err != nil {
		switch err.Error() {
		case "invalid email or password", "missing email or password":
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		case "account is suspended or disabled":
			helpers.RespondWithError(w, http.StatusForbidden, "Your account is suspended")
		case "station registration is pending admin approval":
			helpers.RespondWithError(w, http.StatusForbidden, "Station registration is pending approval")
		case "station registration was rejected":
			helpers.RespondWithError(w, http.StatusForbidden, "Station registration was rejected")
		case "station registration has been revoked":
			helpers.RespondWithError(w, http.StatusForbidden, "Station access has been revoked")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	security.SetAuthCookies(w, resp.AccessToken, resp.RefreshToken)

	_ = helpers.Encode(w, r, http.StatusOK, resp)
}

func resolveRefreshToken(r *http.Request) string {
	if token, ok := security.ExtractRefreshToken(r); ok {
		return token
	}
	req, err := helpers.Decode[RefreshRequest](r)
	if err != nil && !errors.Is(err, io.EOF) {
		return ""
	}
	return strings.TrimSpace(req.RefreshToken)
}

// HandleMe handles GET /api/v1/auth/me.
func HandleMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	resp, err := GetMe(claims.UserId, claims.Role)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, resp)
}

// HandleLogout handles POST /api/v1/auth/logout.
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if token := resolveRefreshToken(r); token != "" {
		_ = Logout(token)
	}
	security.ClearAuthCookies(w)

	_ = helpers.Encode(w, r, http.StatusOK, LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

// HandleRefresh handles POST /api/v1/auth/refresh.
func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	token := resolveRefreshToken(r)
	if token == "" {
		security.ClearAuthCookies(w)
		helpers.RespondWithError(w, http.StatusUnauthorized, "refresh token is required")
		return
	}

	resp, err := RefreshAccessToken(token)
	if err != nil {
		security.ClearAuthCookies(w)
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	security.SetAccessTokenCookie(w, resp.AccessToken)

	_ = helpers.Encode(w, r, http.StatusOK, resp)
}
