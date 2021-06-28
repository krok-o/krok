package models

import (
	"time"
)

// CommandRun is a single run of a command belonging to an event
// including things like, state, event, and created at.
// swagger:model
type CommandRun struct {
	// ID is a generatd identifier.
	//
	// required: true
	ID int `json:"id"`
	// EventID is the ID of the event that this run belongs to.
	//
	// required: true
	EventID int `json:"event_id"`
	// CommandName is the name of the command that is being executed.
	//
	// required: true
	CommandName string `json:"command_name"`
	// Status is the current state of the command run.
	//
	// required: true
	// example: running, failed, success
	Status string `json:"status"`
	// Outcome is any output of the command. Stdout and stderr combined.
	//
	// required: false
	Outcome string `json:"outcome"`
	// CreatedAt is the time when this command run was created.
	//
	// required: true
	CreateAt time.Time `json:"create_at"`
}
