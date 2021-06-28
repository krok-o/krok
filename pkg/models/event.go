package models

import (
	"time"
)

// Event contains details about a platform event, such as
// the repository it belongs to and the event that created it...
// swagger:model
type Event struct {
	// ID is a generated ID.
	//
	// required: true
	ID int `json:"id"`
	// EvenID is the ID of the corresponding event on the given platform. If that cannot be determined
	// an ID is generated.
	//
	// required: true
	EventID string `json:"event_id"`
	// CreatedAt contains the timestamp when this event occurred.
	//
	// required: true
	CreateAt time.Time `json:"create_at"`
	// RepositoryID contains the ID of the repository for which this event occurred.
	//
	// required: true
	RepositoryID int `json:"repository_id"`
	// CommandRuns contains a list of CommandRuns which executed for this event.
	//
	// required: false
	CommandRuns []*CommandRun `json:"command_runs"`
	// Payload defines the information received from the platform for this event.
	//
	// required: true
	Payload string `json:"payload"`
	// VCS is the ID of the platform on which this even occurred.
	//
	// required: true
	VCS int `json:"vcs"`
}
