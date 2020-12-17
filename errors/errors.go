package errors

import "errors"

// NotFound is returned is a resource cannot be found.
var NotFound = errors.New("not found")

// InvalidArgument is sent whenever we encounter an invalid argument.
var InvalidArgument = errors.New("invalid argument")

// QueryError defines an error which occurs when doing database operations.
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string { return e.Query + ": " + e.Err.Error() }
func (e *QueryError) Unwrap() error { return e.Err }
