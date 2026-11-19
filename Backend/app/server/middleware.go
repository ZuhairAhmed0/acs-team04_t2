package server

import (
	"net/http"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/helpers"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/security"
)

// AuthMiddleware validates the JWT and injects claims into the request context.
// Use this to protect routes that require authentication.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, ok := security.ExtractAccessToken(r)
		if !ok {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing or malformed authentication token")
			return
		}
		claims, err := security.ValidateToken(tokenStr)
		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		ctx := security.ContextWithUser(r.Context(), claims)
		next(w, r.WithContext(ctx))
	}
}

// OptionalAuthMiddleware injects claims when a valid JWT is present, but does not block requests without one.
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tokenStr, ok := security.ExtractAccessToken(r); ok {
			if claims, err := security.ValidateToken(tokenStr); err == nil {
				r = r.WithContext(security.ContextWithUser(r.Context(), claims))
			}
		}
		next(w, r)
	}
}

// RequireRole enforces that the authenticated user possesses one of the allowed roles.
func RequireRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := security.UserFromContext(r.Context())
			if !ok || claims == nil {
				helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			allowed := false
			for _, role := range roles {
				if claims.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				helpers.RespondWithError(w, http.StatusForbidden, "Forbidden: insufficient role permissions")
				return
			}

			next(w, r)
		})
	}
}
