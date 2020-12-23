package models

import (
	"time"
)

// User is a user in the Krok system.
type User struct {
	DisplayName string    `json:"display_name,omitempty"`
	Email       string    `json:"email"`
	ID          int       `json:"id"`
	LastLogin   time.Time `json:"last_login,omitempty"`
	APIKeys     []*APIKey `json:"api_keys,omitempty"`
}
