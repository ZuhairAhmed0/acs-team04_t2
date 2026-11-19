/*
	TODO Leaky memory management!!

Change the reteliming strategy from in-memory to Leaky bocket strategy
*/
package security

import (
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/helpers"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// client tracks the rate limiter and last activity time for a single IP.
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit is a per-IP rate limiting middleware.
func RateLimit(limit rate.Limit, burst int) func(http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			mu.Lock()
			c, found := clients[ip]
			if !found {
				c = &client{
					limiter: rate.NewLimiter(limit, burst),
				}
				clients[ip] = c
			}
			c.lastSeen = time.Now()

			if !c.limiter.Allow() {
				mu.Unlock()
				helpers.RespondWithError(w, http.StatusTooManyRequests, "Too many requests. Please try again later.")
				return
			}
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders adds standard security hardening headers to every response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
