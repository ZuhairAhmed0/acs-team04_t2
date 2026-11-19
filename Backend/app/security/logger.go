package security

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func InitLogger(logDir string, defaultWriter io.Writer) (io.Writer, func() error, error) {
	if logDir != "" && logDir != "." {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	dateStr := time.Now().Format("2006-01-02")
	finalLogPath := filepath.Join(logDir, fmt.Sprintf("%s.log", dateStr))

	logFile, err := os.OpenFile(finalLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file %s: %w", finalLogPath, err)
	}

	mw := io.MultiWriter(defaultWriter, logFile)
	log.SetOutput(mw)

	return mw, logFile.Close, nil	
}


type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// LoggingMiddleware logs HTTP requests with timing, client IP, status, and auth user details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		userStr := "anonymous"
		if claims, ok := UserFromContext(r.Context()); ok && claims != nil {
			userStr = claims.UserId + " (" + claims.Role + ")"
		}

		log.Printf(
			"[%s] %s %s | Status: %d | User: %s | IP: %s | Duration: %v",
			r.Method,
			r.URL.Path,
			r.Proto,
			wrapped.statusCode,
			userStr,
			ip,
			duration,
		)
	})
}
