package security

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Standard roles in Fuel Management System.
const (
	RoleUser    = "user"
	RoleStation = "station"
	RoleAdmin   = "admin"
)

// Cookie names for HttpOnly token storage.
const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

// Package-level JWT settings initialized at startup via InitJWT.
var (
	JWTKey          []byte
	AccessDuration  time.Duration
	RefreshDuration time.Duration
	CookieSecure    bool
)

// InitJWT configures the package-level JWT secret key and token durations.
func InitJWT(key []byte, accessDuration, refreshDuration time.Duration) {
	JWTKey = key
	AccessDuration = accessDuration
	RefreshDuration = refreshDuration
}

// SetCookieSecure configures whether auth cookies set the Secure flag.
func SetCookieSecure(secure bool) {
	CookieSecure = secure
}

// Claims holds the JWT payload for authenticated users.
type contextKey string

const UserContextKey contextKey = "user_claims"

type Claims struct {
	UserId      string   `json:"userId"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions,omitempty"`
}

// jwtClaims is an internal struct used to encode/decode the JWT payload.
type jwtClaims struct {
	Claims
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given user using the package-level JWT configuration.
func GenerateToken(userId, role string, permissions []string) (string, error) {
	c := jwtClaims{
		Claims: Claims{
			UserId:      userId,
			Role:        role,
			Permissions: permissions,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(JWTKey)
}

// ValidateToken parses and validates a JWT string using the package-level JWTKey.
func ValidateToken(tokenStr string) (*Claims, error) {
	c := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWTKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	return &c.Claims, nil
}

// ContextWithUser stores claims in the request context.
func ContextWithUser(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, UserContextKey, claims)
}

// UserFromContext retrieves claims from the request context.
func UserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	return claims, ok
}

// GenerateRefreshToken returns a random plaintext token and its SHA-256 hash.
func GenerateRefreshToken() (plaintext string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	plaintext = hex.EncodeToString(b)
	hash = HashToken(plaintext)
	return plaintext, hash, nil
}

// HashToken returns the SHA-256 hex digest of a token string.
func HashToken(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

// ExtractBearerToken parses the Bearer token from the Authorization header.
func ExtractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return strings.TrimSpace(parts[1]), true
}

// ExtractAccessToken reads the JWT from the HttpOnly access_token cookie,
// falling back to the Authorization Bearer header for non-browser clients.
func ExtractAccessToken(r *http.Request) (string, bool) {
	if c, err := r.Cookie(AccessTokenCookieName); err == nil {
		if val := strings.TrimSpace(c.Value); val != "" {
			return val, true
		}
	}
	return ExtractBearerToken(r)
}

// ExtractRefreshToken reads the refresh token from the HttpOnly refresh_token cookie.
func ExtractRefreshToken(r *http.Request) (string, bool) {
	c, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", false
	}
	val := strings.TrimSpace(c.Value)
	return val, val != ""
}

func cookieSameSite() http.SameSite {
	if CookieSecure {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

// SetAccessTokenCookie writes the access JWT into an HttpOnly cookie.
func SetAccessTokenCookie(w http.ResponseWriter, accessToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(AccessDuration.Seconds()),
		Expires:  time.Now().Add(AccessDuration),
		HttpOnly: true,
		Secure:   CookieSecure,
		SameSite: cookieSameSite(),
	})
}

// SetRefreshTokenCookie writes the refresh token into an HttpOnly cookie scoped to auth endpoints.
func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		MaxAge:   int(RefreshDuration.Seconds()),
		Expires:  time.Now().Add(RefreshDuration),
		HttpOnly: true,
		Secure:   CookieSecure,
		SameSite: cookieSameSite(),
	})
}

// SetAuthCookies writes both access and refresh tokens into HttpOnly cookies.
func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	SetAccessTokenCookie(w, accessToken)
	SetRefreshTokenCookie(w, refreshToken)
}

// ClearAuthCookies expires both auth cookies immediately.
func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   CookieSecure,
		SameSite: cookieSameSite(),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   CookieSecure,
		SameSite: cookieSameSite(),
	})
}
