package server

import (
	"net/http"

	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/packages/admin"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/packages/auth"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/packages/station"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/packages/transfers"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/packages/user"
)

// addRoutes registers all API routes on the mux matching the exact Fuel Management contract.
func addRoutes(mux *http.ServeMux) {
	// ─── Authentication Endpoints ───────────────────────────────────────────────
	mux.HandleFunc("POST /api/v1/auth/register/user", auth.HandleRegisterUser)
	mux.HandleFunc("POST /api/v1/auth/register/station", auth.HandleRegisterStation)
	mux.HandleFunc("POST /api/v1/auth/login", auth.HandleLogin)
	mux.HandleFunc("POST /api/v1/auth/refresh", auth.HandleRefresh)
	mux.HandleFunc("GET /api/v1/auth/me", AuthMiddleware(auth.HandleMe))
	mux.HandleFunc("POST /api/v1/auth/logout", auth.HandleLogout)

	// ─── User Endpoints ────────────────────────────────────────────────────────
	mux.HandleFunc("GET /api/v1/user/dashboard", RequireRole("user", "admin")(user.HandleDashboard))
	mux.HandleFunc("GET /api/v1/user/wallet", RequireRole("user", "admin")(user.HandleWallet))
	mux.HandleFunc("GET /api/v1/user/stations", RequireRole("user", "admin")(user.HandleStations))
	mux.HandleFunc("GET /api/v1/user/transactions", RequireRole("user", "admin")(user.HandleTransactions))
	mux.HandleFunc("GET /api/v1/user/transactions/export", RequireRole("user", "admin")(user.HandleExportTransactions))

	// ─── Transfers Endpoints ───────────────────────────────────────────────────
	mux.HandleFunc("POST /api/v1/transfers/topup", RequireRole("user", "admin")(transfers.HandleTopup))
	mux.HandleFunc("POST /api/v1/transfers/allocate", RequireRole("user", "admin")(transfers.HandleAllocate))
	mux.HandleFunc("POST /api/v1/transfers/purchase", RequireRole("user", "admin")(transfers.HandlePurchaseFuel))
	mux.HandleFunc("POST /api/v1/transfers/station-to-station", RequireRole("user", "admin")(transfers.HandleStationToStation))
	mux.HandleFunc("GET /api/v1/transfers", RequireRole("user", "admin")(transfers.HandleGetTransfers))
	mux.HandleFunc("GET /api/v1/transfers/{id}", RequireRole("user", "admin")(transfers.HandleGetTransfer))

	// ─── Station Endpoints ─────────────────────────────────────────────────────
	mux.HandleFunc("GET /api/v1/station/dashboard", RequireRole("station", "admin")(station.HandleDashboard))
	mux.HandleFunc("GET /api/v1/station/wallet", RequireRole("station", "admin")(station.HandleWallet))
	mux.HandleFunc("GET /api/v1/station/transactions", RequireRole("station", "admin")(station.HandleTransactions))
	mux.HandleFunc("GET /api/v1/station/transactions/export", RequireRole("station", "admin")(station.HandleExportTransactions))
	mux.HandleFunc("GET /api/v1/station/branches", RequireRole("station", "admin")(station.HandleGetBranches))
	mux.HandleFunc("POST /api/v1/station/branches", RequireRole("station", "admin")(station.HandleAddBranch))
	mux.HandleFunc("PATCH /api/v1/station/branches/{id}/fuel-status", RequireRole("station", "admin")(station.HandleUpdateFuelStatus))
	mux.HandleFunc("PATCH /api/v1/station/branches/{id}/congestion", RequireRole("station", "admin")(station.HandleUpdateCongestion))
	mux.HandleFunc("POST /api/v1/station/cash-withdrawals", RequireRole("station", "admin")(station.HandleSubmitWithdrawal))
	mux.HandleFunc("GET /api/v1/station/cash-withdrawals", RequireRole("station", "admin")(station.HandleGetWithdrawals))

	// ─── Admin Endpoints ───────────────────────────────────────────────────────
	mux.HandleFunc("GET /api/v1/admin/dashboard", RequireRole("admin")(admin.HandleDashboard))
	mux.HandleFunc("GET /api/v1/admin/users", RequireRole("admin")(admin.HandleGetUsers))
	mux.HandleFunc("PATCH /api/v1/admin/users/{id}/status", RequireRole("admin")(admin.HandleUpdateUserStatus))
	mux.HandleFunc("PATCH /api/v1/admin/users/{id}/role", RequireRole("admin")(admin.HandleUpdateUserRole))
	mux.HandleFunc("GET /api/v1/admin/stations/pending", RequireRole("admin")(admin.HandleGetPendingStations))
	mux.HandleFunc("GET /api/v1/admin/stations", RequireRole("admin")(admin.HandleGetStations))
	mux.HandleFunc("PATCH /api/v1/admin/stations/{id}/approve", RequireRole("admin")(admin.HandleApproveStation))
	mux.HandleFunc("PATCH /api/v1/admin/stations/{id}/reject", RequireRole("admin")(admin.HandleRejectStation))
	mux.HandleFunc("PATCH /api/v1/admin/stations/{id}/revoke", RequireRole("admin")(admin.HandleRevokeStation))
	mux.HandleFunc("DELETE /api/v1/admin/stations/{id}", RequireRole("admin")(admin.HandleDeleteStation))
	mux.HandleFunc("GET /api/v1/admin/cash-withdrawals", RequireRole("admin")(admin.HandleGetWithdrawals))
	mux.HandleFunc("PATCH /api/v1/admin/cash-withdrawals/{id}/approve", RequireRole("admin")(admin.HandleApproveWithdrawal))
	mux.HandleFunc("PATCH /api/v1/admin/cash-withdrawals/{id}/reject", RequireRole("admin")(admin.HandleRejectWithdrawal))
	mux.HandleFunc("PATCH /api/v1/admin/cash-withdrawals/{id}/complete", RequireRole("admin")(admin.HandleCompleteWithdrawal))
}
