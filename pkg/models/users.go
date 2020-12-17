package models

import (
	"time"
)

// User is a user in the Krok system.
type User struct {
	Username    string    `json:"username"`
	ID          string    `json:"id"`
	LastLogin   time.Time `json:"last_login,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
}
