package clip

import (
	"bytes"
	"errors"
)

// NewExitError creates an error with an associated exit code.
func NewExitError(message string, code int) error {
	return exitError{
		code:    code,
		message: message,
	}
}

// exitError is an error with an associated exit code.
type exitError struct {
	code    int
	message string
}

func (e exitError) Error() string { return e.message }
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
func (e usageError) ErrorContext() string {
	b := new(bytes.Buffer)
	_ = writeCommandHelp(b, e.context)
	return b.String()
}
