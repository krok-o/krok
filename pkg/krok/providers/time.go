package providers

import "time"

// Clock wraps time based functions for mocking.
type Clock interface {
	Now() time.Time
}

// KrokClock defines a clock for Krok for testing purposes.
type KrokClock struct{}

// NewClock creates a new KrokClock.
func NewClock() *KrokClock {
	return &KrokClock{}
}

// Now gets the current time.
func (*KrokClock) Now() time.Time {
	return time.Now()
}
