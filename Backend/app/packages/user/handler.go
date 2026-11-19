package user

import (
	"net/http"

	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/helpers"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/security"
)

// HandleDashboard handles GET /api/v1/user/dashboard.
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := ParseUserID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := GetDashboard(userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
	})
}

// HandleWallet handles GET /api/v1/user/wallet.
func HandleWallet(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := ParseUserID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := GetWallet(userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    wallet,
	})
}

// HandleStations handles GET /api/v1/user/stations.
func HandleStations(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	items, err := GetStations(city)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    items,
	})
}

// HandleTransactions handles GET /api/v1/user/transactions.
func HandleTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := ParseUserID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	q := r.URL.Query()
	txType := q.Get("type")
	status := q.Get("status")
	fromDate := q.Get("from_date")
	toDate := q.Get("to_date")

	items, err := GetTransactions(userID, txType, status, fromDate, toDate)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    items,
	})
}

// HandleExportTransactions handles GET /api/v1/user/transactions/export.
func HandleExportTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := ParseUserID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	q := r.URL.Query()
	txType := q.Get("type")
	status := q.Get("status")
	fromDate := q.Get("from_date")
	toDate := q.Get("to_date")

	csvData, err := ExportTransactionsCSV(userID, txType, status, fromDate, toDate)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to export CSV")
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="transactions.csv"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(csvData)
}
