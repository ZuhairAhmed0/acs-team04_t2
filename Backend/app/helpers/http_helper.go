package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// ErrorResponse represents a standardized error format.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RespondWithError writes a standardized JSON error to the response.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Validator is optionally implemented by DTOs that need custom validation
// logic on top of the built-in required-field check.
type Validator interface {
	Validate() error
}

// Encode writes v as JSON to the response with the given status code.
func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
	}

// Decode reads a JSON request body into T.
func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// DecodeValid decodes a JSON request body into T, then:
//  1. Checks every field is non-zero (all fields are required by default).
//     Tag a field with validate:"optional" to skip the check for that field.
//  2. Calls Validate() if T implements the Validator interface, for any
//     extra custom rules beyond the required check.
func DecodeValid[T any](r *http.Request) (T, error) {
	v, err := Decode[T](r)
	if err != nil {
		return v, err
	}
	if err := checkRequired(&v); err != nil {
		return v, err
	}
	if val, ok := any(&v).(Validator); ok {
		if err := val.Validate(); err != nil {
			return v, fmt.Errorf("validate: %w", err)
		}
	}
	return v, nil
}

// checkRequired walks every exported field of a struct.
// All fields are required by default — a field is skipped only when it carries
// the validate:"optional" tag. Returns an error for the first zero-value field.
// Uses only the standard library — no external dependencies.
func checkRequired(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if isOptional(field.Tag.Get("validate")) {
			continue
		}

		if val.Field(i).IsZero() {
			return fmt.Errorf("%s is required", jsonName(field))
		}
	}
	return nil
}

// isOptional reports whether the validate tag includes "optional".
func isOptional(tag string) bool {
	for _, part := range strings.Split(tag, ",") {
		if strings.TrimSpace(part) == "optional" {
			return true
		}
	}
	return false
}

// jsonName returns the JSON field name from the json struct tag,
// falling back to the lowercased Go field name.
func jsonName(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" || tag == "-" {
		return strings.ToLower(f.Name)
	}
	return strings.SplitN(tag, ",", 2)[0]
}

// ResolveIdempotencyKey reads the idempotency key from the Idempotency-Key header
// or falls back to bodyKey, ensuring it is non-empty and within length limits.
func ResolveIdempotencyKey(r *http.Request, bodyKey string) (string, error) {
	key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
	if key == "" {
		key = strings.TrimSpace(bodyKey)
	}
	if key == "" {
		return "", errors.New("idempotency_key is required")
	}
	if len(key) > 120 {
		return "", errors.New("idempotency_key must be at most 120 characters")
	}
	return key, nil
}

// HashIdempotencyPayload returns a deterministic SHA-256 hex digest for request parameters.
func HashIdempotencyPayload(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
