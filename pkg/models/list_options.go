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
	// example: "partialNameOfACommand"
	Name string `json:"name,omitempty"`
	// Only list all entries for a given platform ID.
	//
	// required: false
	// example: 1
	VCS int `json:"vcs,omitempty"`
	// Page defines the current page.
	//
	// required: false
	// example: 0
	Page int `json:"page,omitempty"`
	// PageSize defines the number of items per page.
	//
	// required false
	// example: 10
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
