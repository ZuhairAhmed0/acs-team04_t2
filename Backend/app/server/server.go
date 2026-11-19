package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rs/cors"
	"golang.org/x/time/rate"

	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/config"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/database"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/security"
)

// ServerOptions holds configuration for the HTTP server.
type ServerOptions struct {
	RateLimitRPS   float64
	RateLimitBurst int
	AllowedOrigins []string
}

// NewServer builds and returns the HTTP handler with all middleware applied.
func NewServer(opts ServerOptions) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux)

	var handler http.Handler = mux

	handler = security.RateLimit(rate.Limit(opts.RateLimitRPS), opts.RateLimitBurst)(handler)
	handler = security.SecurityHeaders(handler)
	handler = security.LoggingMiddleware(handler)

	corsOpts := cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Idempotency-Key"},
		AllowCredentials: true,
	}
	if slices.Contains(opts.AllowedOrigins, "*") {
		corsOpts.AllowOriginFunc = func(origin string) bool { return origin != "" }
	} else {
		corsOpts.AllowedOrigins = opts.AllowedOrigins
	}
	c := cors.New(corsOpts)
	handler = c.Handler(handler)

	return handler
}

// Run initializes configuration, databases, storage, and the HTTP server,
// then blocks until the context is cancelled (e.g. SIGINT) and performs a
// graceful shutdown. It is the application entry point called from main.
func Run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Load configuration strictly from system environment with NO fallbacks
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Setup logging to the configured path
	mw, closeLog, err := security.InitLogger(cfg.LogPath, stdout)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer closeLog()
	stdout = mw

	if err := database.InitDB(cfg.DBDsn); err != nil {
		return fmt.Errorf("failed to initialize db at %s: %w", cfg.DBDsn, err)
	}

	security.InitJWT(cfg.JWTSecretKey, cfg.JWTExpiration, cfg.RefreshTokenExpiration)
	security.SetCookieSecure(strings.EqualFold(cfg.AppEnv, "production"))

	// Create server handler using strict configuration
	srv := NewServer(ServerOptions{
		RateLimitRPS:   cfg.RateLimitRPS,
		RateLimitBurst: cfg.RateLimitBurst,
		AllowedOrigins: cfg.CORSAllowedOrigins,
	})

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Port),
		Handler: srv,
	}

	go func() {
		fmt.Fprintf(stdout, "Starting server on port %s (env: %s)\n", cfg.Port, cfg.AppEnv)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, shutdownCancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		} else {
			fmt.Fprintf(stdout, "\nServer gracefully stopped")
		}
	}()
	wg.Wait()
	return nil
}
