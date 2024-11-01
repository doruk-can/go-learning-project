package middleware

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"time"
)

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start timer
			start := time.Now()

			// Create a response writer that captures the status code
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Process the request
			next.ServeHTTP(ww, r)

			// Log the request details
			logger.Info("HTTP request",
				"method", r.Method,
				"uri", r.RequestURI,
				"status", ww.statusCode,
				"duration", time.Since(start).Milliseconds(),
			)
		})
	}
}

// responseWriter is a wrapper to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code for logging
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
