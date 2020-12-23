package errors

import "errors"

// ErrNotFound is returned is a resource cannot be found.
var ErrNotFound = errors.New("not found")

// ErrNoRowsAffected is sent whenever there is a successful query, but no
// rows were affected with the query.
var ErrNoRowsAffected = errors.New("no rows affected")

// ErrAcquireLockFailed signals that the lock for a file name is already taken.
var ErrAcquireLockFailed = errors.New("failed to acquire lock")

// QueryError defines an error which occurs when doing database operations.
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string { return e.Query + ": " + e.Err.Error() }
func (e *QueryError) Unwrap() error { return e.Err }
