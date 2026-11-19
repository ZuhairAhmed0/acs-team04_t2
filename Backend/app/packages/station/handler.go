package station

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/helpers"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/security"
)

// HandleDashboard handles GET /api/v1/station/dashboard.
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := GetDashboard(stationID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
	})
}

// HandleWallet handles GET /api/v1/station/wallet.
func HandleWallet(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := GetWallet(stationID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    wallet,
	})
}

// HandleTransactions handles GET /api/v1/station/transactions.
func HandleTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	q := r.URL.Query()
	txType := q.Get("type")
	status := q.Get("status")
	fromDate := q.Get("from_date")
	toDate := q.Get("to_date")

	txs, err := GetTransactions(stationID, txType, status, fromDate, toDate)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    txs,
	})
}

// HandleExportTransactions handles GET /api/v1/station/transactions/export.
func HandleExportTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	q := r.URL.Query()
	txType := q.Get("type")
	status := q.Get("status")
	fromDate := q.Get("from_date")
	toDate := q.Get("to_date")

	csvData, err := ExportTransactionsCSV(stationID, txType, status, fromDate, toDate)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to export CSV")
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="transactions.csv"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(csvData)
}

// HandleGetBranches handles GET /api/v1/station/branches.
func HandleGetBranches(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	branches, err := GetBranches(stationID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    branches,
	})
}

// HandleAddBranch handles POST /api/v1/station/branches.
func HandleAddBranch(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	req, err := helpers.DecodeValid[AddBranchRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	branch, err := AddBranch(stationID, req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_ = helpers.Encode(w, r, http.StatusCreated, map[string]any{
		"success": true,
		"data":    branch,
	})
}

// HandleUpdateFuelStatus handles PATCH /api/v1/station/branches/{id}/fuel-status.
func HandleUpdateFuelStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	idStr := r.PathValue("id")
	branchID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid branch ID")
		return
	}

	req, err := helpers.DecodeValid[UpdateFuelStatusRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := UpdateFuelStatus(stationID, branchID, req.FuelStatus); err != nil {
		switch err.Error() {
		case "invalid fuel status (must be 'full', 'low', or 'empty')":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		case "branch not found":
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Fuel status updated successfully",
	})
}

// HandleUpdateCongestion handles PATCH /api/v1/station/branches/{id}/congestion.
func HandleUpdateCongestion(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	idStr := r.PathValue("id")
	branchID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid branch ID")
		return
	}

	req, err := helpers.DecodeValid[UpdateCongestionRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := UpdateCongestion(stationID, branchID, req.CongestionPercentage); err != nil {
		switch err.Error() {
		case "congestion_percentage must be between 0 and 100":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		case "branch not found":
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"message": "Congestion percentage updated successfully",
	})
}

// HandleSubmitWithdrawal handles POST /api/v1/station/cash-withdrawals.
func HandleSubmitWithdrawal(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	req, err := helpers.DecodeValid[CashWithdrawalRequest](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := helpers.ResolveIdempotencyKey(r, req.IdempotencyKey)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.IdempotencyKey = key

	result, err := SubmitCashWithdrawal(stationID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			helpers.RespondWithError(w, http.StatusConflict, err.Error())
		case err.Error() == "amount must be greater than zero":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		case err.Error() == "insufficient available balance":
			helpers.RespondWithError(w, http.StatusPaymentRequired, err.Error())
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    result,
	})
}

// HandleGetWithdrawals handles GET /api/v1/station/cash-withdrawals.
func HandleGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stationID, err := ParseStationID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	list, err := GetCashWithdrawals(stationID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    list,
	})
}
