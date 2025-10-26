package errors

import (
	"encoding/json"
	"net/http"

	"github.com/F1sssss/Perfect_Trade/internal/shared/logger"
)

// ErrorResponse is the JSON error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// WriteError writes an error as JSON response
func WriteError(w http.ResponseWriter, r *http.Request, err error, log logger.Logger) {
	// Determine status code and response based on error type
	status, code, message := classifyError(err)

	// Log the error with full details
	log.Error("request failed",
		logger.String("path", r.URL.Path),
		logger.String("method", r.Method),
		logger.Int("status", status),
		logger.String("code", code),
		logger.Error(err),
	)

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// classifyError determines HTTP status code based on error type
func classifyError(err error) (status int, code string, message string) {
	switch {
	case Is(err, ErrValidation):
		return http.StatusBadRequest, "VALIDATION_ERROR", err.Error()
	case Is(err, ErrNotFound):
		return http.StatusNotFound, "NOT_FOUND", err.Error()
	case Is(err, ErrAlreadyExists):
		return http.StatusConflict, "CONFLICT", err.Error()
	case Is(err, ErrInvalidInput):
		return http.StatusBadRequest, "INVALID_INPUT", err.Error()
	case Is(err, ErrBusinessRule):
		return http.StatusBadRequest, "BUSINESS_RULE_VIOLATION", err.Error()
	case Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required"
	case Is(err, ErrForbidden):
		return http.StatusForbidden, "FORBIDDEN", "Access denied"
	default:
		// Unknown error - treat as internal server error
		// Don't leak internal details to client
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred"
	}
}
