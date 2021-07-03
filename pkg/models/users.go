package models

import (
	"time"
)

// User is a user in the Krok system.
// swagger:model
type User struct {
	// DisplayName is the name of the user.
	//
	// required: false
	DisplayName string `json:"display_name,omitempty"`
	// Email of the user.
	//
	// required: true
	Email string `json:"email"`
	// ID of the user. This is auto-generated.
	//
	// required: true
	ID int `json:"id"`
	// LastLogin contains the timestamp of the last successful login of this user.
	//
	// required: true
	LastLogin time.Time `json:"last_login,omitempty"`
	// APIKeys contains generated api access keys for this user.
	//
	// required: false
	APIKeys []*APIKey `json:"api_keys,omitempty"`
}
