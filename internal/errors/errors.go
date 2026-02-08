package errors

import (
	"fmt"
	"strings"
)

// Error represents a PlatoSL error with context
type Error struct {
	Type       ErrorType
	Message    string
	File       string
	Line       int
	Column     int
	Suggestion string
	Cause      error
}

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeConfig        ErrorType = "config"
	ErrorTypeGeneration    ErrorType = "generation"
	ErrorTypeFileSystem    ErrorType = "filesystem"
	ErrorTypeInternal      ErrorType = "internal"
)

// New creates a new error
func New(typ ErrorType, msg string) *Error {
	return &Error{
		Type:    typ,
		Message: msg,
	}
}

// Newf creates a new error with formatted message
func Newf(typ ErrorType, format string, args ...interface{}) *Error {
	return &Error{
		Type:    typ,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an existing error
func Wrap(typ ErrorType, cause error, msg string) *Error {
	return &Error{
		Type:    typ,
		Message: msg,
		Cause:   cause,
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(typ ErrorType, cause error, format string, args ...interface{}) *Error {
	return &Error{
		Type:    typ,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

// WithLocation adds location information
func (e *Error) WithLocation(file string, line, column int) *Error {
	e.File = file
	e.Line = line
	e.Column = column
	return e
}

// WithSuggestion adds a suggestion
func (e *Error) WithSuggestion(suggestion string) *Error {
	e.Suggestion = suggestion
	return e
}

// Error implements the error interface
func (e *Error) Error() string {
	var b strings.Builder

	// Type prefix
	if e.Type != "" {
		fmt.Fprintf(&b, "[%s] ", e.Type)
	}

	// Location
	if e.File != "" {
		fmt.Fprintf(&b, "%s", e.File)
		if e.Line > 0 {
			fmt.Fprintf(&b, ":%d", e.Line)
			if e.Column > 0 {
				fmt.Fprintf(&b, ":%d", e.Column)
			}
		}
		b.WriteString(": ")
	}

	// Message
	b.WriteString(e.Message)

	// Cause
	if e.Cause != nil {
		fmt.Fprintf(&b, ": %v", e.Cause)
	}

	return b.String()
}

// Format formats the error for user display
func (e *Error) Format() string {
	var b strings.Builder

	// Location header
	if e.File != "" {
		fmt.Fprintf(&b, "✗ %s", e.File)
		if e.Line > 0 {
			fmt.Fprintf(&b, ":%d", e.Line)
			if e.Column > 0 {
				fmt.Fprintf(&b, ":%d", e.Column)
			}
		}
		b.WriteString(": ")
		b.WriteString(e.Message)
	} else {
		fmt.Fprintf(&b, "✗ %s", e.Message)
	}

	// Cause details
	if e.Cause != nil {
		fmt.Fprintf(&b, "\n\n  Error: %v", e.Cause)
	}

	// Suggestion
	if e.Suggestion != "" {
		fmt.Fprintf(&b, "\n\n  Suggestion: %s", e.Suggestion)
	}

	return b.String()
}

// FormatMultiple formats multiple errors
func FormatMultiple(errs []*Error) string {
	if len(errs) == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Found %d error(s):\n\n", len(errs))

	for i, err := range errs {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(err.Format())
		b.WriteString("\n")
	}

	return b.String()
}
