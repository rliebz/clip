package clip

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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

// getExitCode gets an exit code if it exists, or returns 1.
func getExitCode(err error) int {
	var exitErr exitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return 1
}

// printError prints an error with contextual information.
func printError(l *log.Logger, err error) {
	l.Printf("Error: %s\n", err)

	if ectx, ok := err.(errorContext); ok {
		l.Println()
		l.Print(ectx.ErrorContext())
	}
}

// errorContext includes a context to be printed after an error.
type errorContext interface {
	ErrorContext() string
}

// newUsageError creates an error which causes help to be printed.
func newUsageError(ctx *Context, message string) usageError {
	return usageError{
		context: ctx,
		message: message,
	}
}

// newUsageErrorf creates an error which causes help to be printed.
func newUsageErrorf(ctx *Context, format string, a ...any) usageError {
	return usageError{
		context: ctx,
		message: fmt.Sprintf(format, a...),
	}
}

// usageError is an error caused by incorrect usage.
type usageError struct {
	context *Context
	message string
}

func (e usageError) Error() string { return e.message }
func (e usageError) ErrorContext() string {
	b := new(bytes.Buffer)
	_ = writeCommandHelp(b, e.context)
	return b.String()
}
