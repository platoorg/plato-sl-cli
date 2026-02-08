package cue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
)

// ValidationResult holds the result of a validation
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// ValidationError represents a single validation error with context
type ValidationError struct {
	File       string
	Line       int
	Column     int
	Path       string
	Message    string
	Suggestion string
}

// Validator validates CUE values
type Validator struct {
	strict bool
}

// NewValidator creates a new validator
func NewValidator(strict bool) *Validator {
	return &Validator{strict: strict}
}

// Validate validates a CUE value
func (v *Validator) Validate(val cue.Value) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	// Check for errors in the value
	if err := val.Err(); err != nil {
		result.Valid = false
		result.Errors = v.parseErrors(err)
		return result
	}

	// Validate the value (concrete check)
	if err := val.Validate(cue.Concrete(v.strict)); err != nil {
		result.Valid = false
		result.Errors = v.parseErrors(err)
		return result
	}

	return result
}

// parseErrors converts CUE errors to ValidationErrors
func (v *Validator) parseErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	// Use CUE's error formatting
	for _, e := range errors.Errors(err) {
		pos := e.Position()

		validationErrors = append(validationErrors, ValidationError{
			File:       pos.Filename(),
			Line:       pos.Line(),
			Column:     pos.Column(),
			Path:       extractPath(e),
			Message:    cleanMessage(e.Error()),
			Suggestion: generateSuggestion(e.Error()),
		})
	}

	return validationErrors
}

// extractPath extracts the field path from an error
func extractPath(err errors.Error) string {
	// Try to extract path from error message
	msg := err.Error()
	if idx := strings.Index(msg, ":"); idx > 0 {
		path := strings.TrimSpace(msg[:idx])
		if !strings.Contains(path, " ") {
			return path
		}
	}
	return ""
}

// cleanMessage cleans up the error message
func cleanMessage(msg string) string {
	// Remove file:line:col prefix if present
	if idx := strings.Index(msg, ": "); idx > 0 {
		after := msg[idx+2:]
		if strings.Contains(msg[:idx], ":") {
			return after
		}
	}
	return msg
}

// generateSuggestion generates a helpful suggestion based on the error
func generateSuggestion(msg string) string {
	lower := strings.ToLower(msg)

	if strings.Contains(lower, "concrete") {
		return "Ensure all fields have concrete values (no unresolved references)"
	}

	if strings.Contains(lower, "conflicting") || strings.Contains(lower, "conflict") {
		return "Check for duplicate or contradicting field definitions"
	}

	if strings.Contains(lower, "incomplete") {
		return "Some required fields may be missing or undefined"
	}

	if strings.Contains(lower, "reference") {
		return "Check that all referenced fields and definitions exist"
	}

	if strings.Contains(lower, "cannot use") {
		return "Type mismatch - check that values match their expected types"
	}

	return ""
}

// FormatError formats a validation error for display
func FormatError(err ValidationError) string {
	var b strings.Builder

	// File:line:col
	if err.File != "" {
		fmt.Fprintf(&b, "%s:%d:%d: ", err.File, err.Line, err.Column)
	}

	// Path
	if err.Path != "" {
		fmt.Fprintf(&b, "field '%s': ", err.Path)
	}

	// Message
	b.WriteString(err.Message)

	// Suggestion
	if err.Suggestion != "" {
		fmt.Fprintf(&b, "\n\n  Suggestion: %s", err.Suggestion)
	}

	return b.String()
}
