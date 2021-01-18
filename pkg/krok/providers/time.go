package providers

import "time"

// Clock wraps time based functions for mocking.
type Clock interface {
	Now() time.Time
}

type clock struct{}

// NewClock creates a new clock.
func NewClock() *clock {
	return &clock{}
}

// Now gets the current time.
func (_ *clock) Now() time.Time {
	return time.Now()
}
