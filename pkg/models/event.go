package models

import (
	"time"
)

// Event contains details about a platform event, such as
// the repository it belongs to and the event that created it...
type Event struct {
	ID           int           `json:"id"`
	EventID      string        `json:"event_id"`
	CreateAt     time.Time     `json:"create_at"`
	RepositoryID int           `json:"repository_id"`
	CommandRuns  []*CommandRun `json:"command_runs"`
	Payload      string        `json:"payload"`
	VCS          int           `json:"vcs"`
}
