package models

import (
	"time"
)

// ListOptions provides options for List operations with additional filters.
type ListOptions struct {
	Name string `json:"name,omitempty"`
	// List all entries for Git, Gitea...
	VCS int `json:"vcs,omitempty"`
	// Current Page
	Page int `json:"page,omitempty"`
	// Items per Page
	PageSize int `json:"page_size,omitempty"`
	// Starting Date
	StartingDate *time.Time `json:"starting_date,omitempty"`
	// Ending Date
	EndDate *time.Time `json:"end_date,omitempty"`
}
