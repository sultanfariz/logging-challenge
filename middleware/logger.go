package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware wraps an http.Handler with logging functionality
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Initialize logger with request ID
		logger := log.With().Str("request_id", uuid.New().String()).Logger()
		ctx := logger.WithContext(r.Context())

		// Create a new response writer that can capture the status code
		lw := NewLoggingResponseWriter(w)

		// Call the next handler with the logger-enriched context
		next.ServeHTTP(lw, r.WithContext(ctx))

		// Log the request completion
		logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", lw.statusCode).
			Msg("request completed")
	})
}

// LoggingResponseWriter wraps http.ResponseWriter to capture the status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter creates a new LoggingResponseWriter
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (lw *LoggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

// LogError centralizes error logging
func LogError(ctx context.Context, function string, err error, msg string) {
	log.Ctx(ctx).Error().
		Str("function", function).
		Err(err).
		Msg(msg)
}
