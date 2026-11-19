package admin

import (
	"net/http"
	"strconv"

	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/helpers"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/security"
)

// HandleDashboard handles GET /api/v1/admin/dashboard.
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := GetDashboard()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
	})
}

// HandleGetUsers handles GET /api/v1/admin/users.
func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := GetUsers()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    users,
	})
}

// HandleUpdateUserStatus handles PATCH /api/v1/admin/users/{id}/status.
func HandleUpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req, err := helpers.DecodeValid[UpdateUserStatusRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ChangeUserStatus(id, req.Status); err != nil {
		switch err.Error() {
		case "invalid status (must be 'active' or 'suspended')":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		case "user not found":
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "User status updated successfully",
	})
}

// HandleUpdateUserRole handles PATCH /api/v1/admin/users/{id}/role.
func HandleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req, err := helpers.DecodeValid[UpdateUserRoleRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ChangeUserRole(id, req.Role); err != nil {
		switch err.Error() {
		case "invalid user role (must be 'user' or 'admin')":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		case "user not found":
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "User role updated successfully",
	})
}

// HandleGetPendingStations handles GET /api/v1/admin/stations/pending.
func HandleGetPendingStations(w http.ResponseWriter, r *http.Request) {
	stations, err := ListPendingStations()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    stations,
	})
}

// HandleGetStations handles GET /api/v1/admin/stations.
func HandleGetStations(w http.ResponseWriter, r *http.Request) {
	stations, err := ListStations()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    stations,
	})
}

// HandleApproveStation handles PATCH /api/v1/admin/stations/{id}/approve.
func HandleApproveStation(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	adminID, err := ParseAdminID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid station ID")
		return
	}

	if err := ApproveStationRegistration(id, adminID); err != nil {
		if err.Error() == "station not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Station approved successfully",
	})
}

// HandleRejectStation handles PATCH /api/v1/admin/stations/{id}/reject.
func HandleRejectStation(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	adminID, err := ParseAdminID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid station ID")
		return
	}

	req, _ := helpers.DecodeValid[RejectStationRequest](r)

	if err := RejectStationRegistration(id, adminID, req.RejectionReason); err != nil {
		if err.Error() == "station not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Station rejected",
	})
}

// HandleRevokeStation handles PATCH /api/v1/admin/stations/{id}/revoke.
func HandleRevokeStation(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid station ID")
		return
	}

	if err := RevokeStationRegistration(id); err != nil {
		if err.Error() == "station not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Station approval revoked",
	})
}

// HandleDeleteStation handles DELETE /api/v1/admin/stations/{id}.
func HandleDeleteStation(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid station ID")
		return
	}

	if err := RemoveStation(id); err != nil {
		if err.Error() == "station not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Station deleted successfully",
	})
}

// HandleGetWithdrawals handles GET /api/v1/admin/cash-withdrawals.
func HandleGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	list, err := ListWithdrawals()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    list,
	})
}

// HandleApproveWithdrawal handles PATCH /api/v1/admin/cash-withdrawals/{id}/approve.
func HandleApproveWithdrawal(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	adminID, err := ParseAdminID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid withdrawal ID")
		return
	}

	if err := ApproveWithdrawal(id, adminID); err != nil {
		if err.Error() == "cash withdrawal request not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Cash withdrawal approved",
	})
}

// HandleRejectWithdrawal handles PATCH /api/v1/admin/cash-withdrawals/{id}/reject.
func HandleRejectWithdrawal(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	adminID, err := ParseAdminID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid withdrawal ID")
		return
	}

	req, _ := helpers.DecodeValid[RejectWithdrawalRequest](r)

	if err := RejectWithdrawal(id, adminID, req.RejectionReason); err != nil {
		if err.Error() == "cash withdrawal request not found" {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Cash withdrawal rejected",
	})
}

// HandleCompleteWithdrawal handles PATCH /api/v1/admin/cash-withdrawals/{id}/complete.
func HandleCompleteWithdrawal(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid withdrawal ID")
		return
	}

	req, _ := helpers.DecodeValid[CompleteWithdrawalRequest](r)

	if err := CompleteWithdrawal(id, req.ExternalReference); err != nil {
		switch err.Error() {
		case "cash withdrawal request not found", "station not found":
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		case "cash withdrawal request is already settled", "cash withdrawal request must be approved before completion":
			helpers.RespondWithError(w, http.StatusConflict, err.Error())
		case "insufficient available balance":
			helpers.RespondWithError(w, http.StatusPaymentRequired, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"id":     id,
			"status": "completed",
		},
	})
}
