package models

import (
	"time"
)

// CommandRun is a single run of a command belonging to an event
// including things like, state, event, and created at.
type CommandRun struct {
	ID       int       `json:"id"`
	EventID  int       `json:"event_id"`
	Status   string    `json:"status"`
	Outcome  string    `json:"outcome"`
	CreateAt time.Time `json:"create_at"`
}
