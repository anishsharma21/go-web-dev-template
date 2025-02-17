package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal/auth"
	"github.com/google/uuid"
)

// Helper responseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// CustomLogHandler for adding request_id and user_id fields to the log record
type CustomLogHandler struct {
	slog.Handler
}

// Handle adds request_id and user_id fields to the log record if they exist in the context
func (h *CustomLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		r.AddAttrs(slog.String("request_id", requestID))
	}
	if userID, ok := ctx.Value("user_id").(int); ok {
		r.AddAttrs(slog.Int("user_id", userID))
	}
	return h.Handler.Handle(ctx, r)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a request ID for traceability
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request_id", requestID)

		// get email from JWT claims if present
		tokenString := r.Header.Get("Authorization")
		if tokenString != "" {
			id, err := auth.ParseIdFromToken(tokenString)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to parse token claims", "error", err)
			}

			ctx = context.WithValue(ctx, "user_id", id)
		}

		// Capture response status code using a response writer wrapper
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r.WithContext(ctx))

		duration := float64(time.Since(start))

		// Log the request and response information
		slog.InfoContext(ctx, fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, rw.statusCode),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status_code", rw.statusCode),
			slog.String("response_time", fmt.Sprintf("%.2fms", duration/1_000_000.0)),
			slog.String("client_ip", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}
