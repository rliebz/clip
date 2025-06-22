package clip

import (
	"bytes"
	"errors"
	"fmt"
)

// exitError is an error with an associated exit code.
type exitError struct {
	code int
	err  error
}

// NewExitError creates an error with an associated exit code.
func NewExitError(code int, message string) error {
	return exitError{
		code: code,
		err:  errors.New(message),
	}
}

// NewExitErrorf creates an error with an associated exit code using
// [fmt.Errorf] semantics.
func NewExitErrorf(code int, format string, a ...any) error {
	return exitError{
		code: code,
		err:  fmt.Errorf(format, a...),
	}
}

// WithExitCode wraps an existing error with an associated exit code.
func WithExitCode(code int, err error) error {
	return exitError{
		code: code,
		err:  err,
	}
}

func (e exitError) Error() string { return e.err.Error() }
func (e exitError) Unwrap() error { return e.err }
func (e exitError) ExitCode() int { return e.code }

// exitCode gets an exit code if it exists, or returns 1 for non-nil errors.
func exitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr interface{ ExitCode() int }
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

// errorContext includes a context to be printed after an error.
type errorContext interface {
	ErrorContext() string
}

// newUsageError creates an error which causes help to be printed.
func newUsageError(ctx *Context, err error) usageError {
	return usageError{
		context: ctx,
		err:     err,
	}
}

// usageError is an error caused by incorrect usage.
type usageError struct {
	context *Context
	err     error
}

func (e usageError) Error() string { return e.err.Error() }
func (e usageError) Unwrap() error { return e.err }
func (e usageError) ExitCode() int { return 2 }
func (e usageError) ErrorContext() string {
	b := new(bytes.Buffer)
	_ = WriteHelp(b, e.context)
	return b.String()
}
