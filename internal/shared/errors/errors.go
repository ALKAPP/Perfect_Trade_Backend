package errors

import (
	"errors"
	"fmt"
)

var (
	// Domain errors (400 Bad Request)
	ErrValidation    = errors.New("validation failed")
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrBusinessRule  = errors.New("business rule violation")

	// Infrastructure errors (500 Internal Server Error)
	ErrDatabase    = errors.New("database error")
	ErrTransaction = errors.New("transaction error")
	ErrExternal    = errors.New("external service error")
	ErrInternal    = errors.New("internal error")

	// Auth errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// Wrap adds context to an error
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf adds formatted context to an error
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// Is checks if err is or wraps target
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As checks if err can be assigned to target
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
