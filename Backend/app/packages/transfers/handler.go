package transfers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"wa/app/helpers"
	"wa/app/security"
)

func decodeTransferRequest[T any](w http.ResponseWriter, r *http.Request) (uint64, T, string, bool) {
	var zero T
	claims, ok := security.UserFromContext(r.Context())
	if !ok || claims == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return 0, zero, "", false
	}

	userID, err := ParseUserID(claims.UserId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return 0, zero, "", false
	}

	req, err := helpers.DecodeValid[T](r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return 0, zero, "", false
	}

	return userID, req, claims.UserId, true
}

// HandleTopup handles POST /api/v1/transfers/topup.
func HandleTopup(w http.ResponseWriter, r *http.Request) {
	userID, req, _, ok := decodeTransferRequest[TopupRequest](w, r)
	if !ok {
		return
	}

	key, err := helpers.ResolveIdempotencyKey(r, req.IdempotencyKey)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.IdempotencyKey = key

	result, err := Topup(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			helpers.RespondWithError(w, http.StatusConflict, err.Error())
		case err.Error() == "amount must be greater than zero":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
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

// HandleAllocate handles POST /api/v1/transfers/allocate.
func HandleAllocate(w http.ResponseWriter, r *http.Request) {
	userID, req, _, ok := decodeTransferRequest[AllocateRequest](w, r)
	if !ok {
		return
	}

	key, err := helpers.ResolveIdempotencyKey(r, req.IdempotencyKey)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.IdempotencyKey = key

	result, err := Allocate(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			helpers.RespondWithError(w, http.StatusConflict, err.Error())
		case err.Error() == "amount must be greater than zero" || err.Error() == "station not found":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		case err.Error() == "insufficient available balance":
			helpers.RespondWithError(w, http.StatusPaymentRequired, err.Error())
		case err.Error() == "station account is not active":
			helpers.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
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

// HandlePurchaseFuel handles POST /api/v1/transfers/purchase.
func HandlePurchaseFuel(w http.ResponseWriter, r *http.Request) {
	userID, req, _, ok := decodeTransferRequest[FuelPurchaseRequest](w, r)
	if !ok {
		return
	}

	key, err := helpers.ResolveIdempotencyKey(r, req.IdempotencyKey)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.IdempotencyKey = key

	result, err := PurchaseFuel(userID, req)
	if err != nil {
		msg := err.Error()
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			helpers.RespondWithError(w, http.StatusConflict, msg)
		case msg == "amount must be greater than zero" ||
			msg == "amount cannot be negative" ||
			msg == "loyalty_points_used or amount must be greater than zero" ||
			msg == "branch not found" ||
			msg == "invalid payment source" ||
			msg == "loyalty_points_used cannot be negative" ||
			msg == "loyalty_points_used does not match purchase amount" ||
			msg == "loyalty discount cannot exceed purchase amount" ||
			msg == "loyalty points are only available for user accounts" ||
			strings.HasPrefix(msg, "loyalty_points_used must be a multiple of") ||
			strings.HasPrefix(msg, "amount must be a multiple of"):
			helpers.RespondWithError(w, http.StatusBadRequest, msg)
		case msg == "insufficient available balance" ||
			msg == "insufficient allocated balance" ||
			msg == "insufficient station balance" ||
			msg == "insufficient loyalty points":
			helpers.RespondWithError(w, http.StatusPaymentRequired, msg)
		case msg == "selected branch is out of fuel":
			helpers.RespondWithError(w, http.StatusConflict, msg)
		case msg == "station not found" || msg == "station account is not active":
			helpers.RespondWithError(w, http.StatusUnprocessableEntity, msg)
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

// HandleStationToStation handles POST /api/v1/transfers/station-to-station.
func HandleStationToStation(w http.ResponseWriter, r *http.Request) {
	userID, req, _, ok := decodeTransferRequest[StationTransferRequest](w, r)
	if !ok {
		return
	}

	key, err := helpers.ResolveIdempotencyKey(r, req.IdempotencyKey)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.IdempotencyKey = key

	result, err := StationToStation(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			helpers.RespondWithError(w, http.StatusConflict, err.Error())
		case err.Error() == "amount must be greater than zero" || err.Error() == "source and destination stations must be different":
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		case err.Error() == "insufficient allocated balance" || err.Error() == "insufficient available balance":
			helpers.RespondWithError(w, http.StatusPaymentRequired, err.Error())
		case err.Error() == "station not found" || err.Error() == "station account is not active":
			helpers.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
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

// HandleGetTransfers handles GET /api/v1/transfers.
func HandleGetTransfers(w http.ResponseWriter, r *http.Request) {
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

	transfers, err := GetTransfers(userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    transfers,
	})
}

// HandleGetTransfer handles GET /api/v1/transfers/{id}.
func HandleGetTransfer(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid transfer ID")
		return
	}

	detail, err := GetTransfer(id, userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Transfer not found")
		return
	}

	_ = helpers.Encode(w, r, http.StatusOK, map[string]any{
		"success": true,
		"data":    detail,
	})
}
