package config

import (
	"fmt"
	"strings"
)

// Validator validates configuration values
type Validator struct {
	errors []string
}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
	return &Validator{
		errors: []string{},
	}
}

// Required checks if a value is not empty
func (v *Validator) Required(key, value string) {
	if value == "" {
		v.errors = append(v.errors, fmt.Sprintf("%s is required", key))
	}
}

// OneOf checks if a value is one of the allowed values
func (v *Validator) OneOf(key, value string, allowed []string) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	v.errors = append(v.errors, fmt.Sprintf("%s must be one of: %s", key, strings.Join(allowed, ", ")))
}

// Range checks if an integer is within a range
func (v *Validator) Range(key string, value, min, max int) {
	if value < min || value > max {
		v.errors = append(v.errors, fmt.Sprintf("%s must be between %d and %d", key, min, max))
	}
}

// Min checks if an integer is at least a minimum value
func (v *Validator) Min(key string, value, min int) {
	if value < min {
		v.errors = append(v.errors, fmt.Sprintf("%s must be at least %d", key, min))
	}
}

// MinLength checks if a string has a minimum length
func (v *Validator) MinLength(key, value string, min int) {
	if len(value) < min {
		v.errors = append(v.errors, fmt.Sprintf("%s must be at least %d characters", key, min))
	}
}

// Error returns all validation errors as a single error
func (v *Validator) Error() error {
	if len(v.errors) == 0 {
		return nil
	}
	return fmt.Errorf("configuration validation failed:\n  - %s", strings.Join(v.errors, "\n  - "))
}
