package models

import (
	"time"
)

// ListOptions provides options for List operations with additional filters.
// swagger:model
type ListOptions struct {
	// Name of the context for which this option is used.
	//
	// required: false
	Name string `json:"name,omitempty"`
	// Only list all entries for a given platform ID.
	//
	// required: false
	VCS int `json:"vcs,omitempty"`
	// Current Page
	//
	// required: false
	Page int `json:"page,omitempty"`
	// Items per Page
	//
	// required false
	PageSize int `json:"page_size,omitempty"`
	// StartingDate defines a date of start to look for events. Inclusive.
	//
	// required: false
	// example: 2021-02-02
	StartingDate *time.Time `json:"starting_date,omitempty"`
	// EndDate defines a date of end to look for events. Not Inclusive.
	//
	// required: false
	// example: 2021-02-03
	EndDate *time.Time `json:"end_date,omitempty"`
}
