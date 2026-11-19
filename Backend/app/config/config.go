package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime application configuration without hardcoded fallbacks.
type Config struct {
	Port                   string
	AppEnv                 string
	JWTSecretKey           []byte
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration
	DBDsn                  string
	LogPath                string
	RateLimitRPS           float64
	RateLimitBurst         int
	CORSAllowedOrigins     []string
}

// Load reads configuration from env files and system environment.
// It strictly validates that all required parameters are provided with no fallbacks.
func Load(envFiles ...string) (*Config, error) {
	if len(envFiles) == 0 {
		envFiles = []string{".env", "../.env"}
	}
	var logEnvDir string
	for _, envFile := range envFiles {
		hadLogPath := strings.TrimSpace(os.Getenv("LOG_PATH")) != ""
		if err := godotenv.Load(envFile); err == nil && !hadLogPath && strings.TrimSpace(os.Getenv("LOG_PATH")) != "" {
			logEnvDir = filepath.Dir(envFile)
		}
	}

	cfg := &Config{}

	// Port
	cfg.Port = os.Getenv("PORT")
	if strings.TrimSpace(cfg.Port) == "" {
		return nil, errors.New("environment variable PORT is required")
	}

	// AppEnv
	cfg.AppEnv = os.Getenv("APP_ENV")
	if strings.TrimSpace(cfg.AppEnv) == "" {
		return nil, errors.New("environment variable APP_ENV is required")
	}

	// JWTSecretKey
	secret := os.Getenv("JWT_SECRET_KEY")
	if strings.TrimSpace(secret) == "" {
		return nil, errors.New("environment variable JWT_SECRET_KEY is required")
	}
	cfg.JWTSecretKey = []byte(secret)

	// JWTExpiration
	expStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if strings.TrimSpace(expStr) == "" {
		return nil, errors.New("environment variable JWT_EXPIRATION_HOURS is required")
	}
	expHours, err := strconv.Atoi(expStr)
	if err != nil || expHours <= 0 {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS value: %q", expStr)
	}
	cfg.JWTExpiration = time.Duration(expHours) * time.Hour

	// RefreshTokenExpirationDays
	refreshStr := os.Getenv("REFRESH_TOKEN_EXPIRATION_DAYS")
	if strings.TrimSpace(refreshStr) == "" {
		return nil, errors.New("environment variable REFRESH_TOKEN_EXPIRATION_DAYS is required")
	}
	refreshDays, err := strconv.Atoi(refreshStr)
	if err != nil || refreshDays <= 0 {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRATION_DAYS value: %q", refreshStr)
	}
	cfg.RefreshTokenExpiration = time.Duration(refreshDays) * 24 * time.Hour

	// DBDsn
	if dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); dbURL != "" {
		cfg.DBDsn = dbURL
	} else if os.Getenv("POSTGRES_HOST") != "" {
		sslmode := strings.TrimSpace(os.Getenv("POSTGRES_SSLMODE"))
		if sslmode == "" {
			sslmode = "disable"
		}
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USERNAME"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DATABASE"),
			sslmode,
		)
		if tz := strings.TrimSpace(os.Getenv("POSTGRES_TIMEZONE")); tz != "" {
			dsn += " TimeZone=" + tz
		}
		cfg.DBDsn = dsn
	} else {
		return nil, errors.New("database configuration required: provide DATABASE_URL or POSTGRES_HOST/PORT/USERNAME/PASSWORD/DATABASE")
	}

	// LogPath (resolve relative paths from the loaded .env directory so logs stay in root/logs)
	cfg.LogPath = strings.TrimSpace(os.Getenv("LOG_PATH"))
	if cfg.LogPath == "" {
		cfg.LogPath = "logs"
	}
	if logEnvDir != "" && !filepath.IsAbs(cfg.LogPath) {
		if absLogPath, err := filepath.Abs(filepath.Join(logEnvDir, cfg.LogPath)); err == nil {
			cfg.LogPath = absLogPath
		}
	}

	// RateLimitRPS
	rpsStr := os.Getenv("RATE_LIMIT_RPS")
	if strings.TrimSpace(rpsStr) == "" {
		return nil, errors.New("environment variable RATE_LIMIT_RPS is required")
	}
	rps, err := strconv.ParseFloat(rpsStr, 64)
	if err != nil || rps <= 0 {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPS value: %q", rpsStr)
	}
	cfg.RateLimitRPS = rps

	// RateLimitBurst
	burstStr := os.Getenv("RATE_LIMIT_BURST")
	if strings.TrimSpace(burstStr) == "" {
		return nil, errors.New("environment variable RATE_LIMIT_BURST is required")
	}
	burst, err := strconv.Atoi(burstStr)
	if err != nil || burst <= 0 {
		return nil, fmt.Errorf("invalid RATE_LIMIT_BURST value: %q", burstStr)
	}
	cfg.RateLimitBurst = burst

	// CORSAllowedOrigins
	originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	if strings.TrimSpace(originsStr) == "" {
		return nil, errors.New("environment variable CORS_ALLOWED_ORIGINS is required")
	}
	var origins []string
	for _, o := range strings.Split(originsStr, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		return nil, errors.New("CORS_ALLOWED_ORIGINS must contain at least one origin")
	}
	cfg.CORSAllowedOrigins = origins

	return cfg, nil
}
