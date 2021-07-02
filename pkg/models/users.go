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
	// Token of the user for API logins using krokctl or curl. This is only displayed for new users.
	// Cannot be retrieved ever again. Regenerate if forgotten or misplaced.
	//
	// required: false
	Token *string `json:"-"`
}

// NewUser is a new user in the Krok system. Specifically this exposes the token and should only be used when creating
// a user for the first time.
// swagger:model
type NewUser struct {
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
	// Token is displayed once for new users. Then never again.
	//
	// required: true
	Token *string `json:"token"`
}
