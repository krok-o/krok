package errors

import "errors"

// ErrNotFound is returned is a resource cannot be found.
var ErrNotFound = errors.New("not found")

// ErrNoRowsAffected is sent whenever there is a successful query, but no
// rows were affected with the query.
var ErrNoRowsAffected = errors.New("no rows affected")

// QueryError defines an error which occurs when doing database operations.
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string { return e.Query + ": " + e.Err.Error() }

// Unwrap unwraps the query error and returns the internal error.
func (e *QueryError) Unwrap() error { return e.Err }

// Message represents an error message.
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// APIError wraps a message and a code into a struct for JSON parsing.
func APIError(m string, code int, err error) Message {
	if err == nil {
		err = errors.New("unexpected error")
	}
	return Message{
		Code:    code,
		Message: m,
		Error:   err.Error(),
	}
}
