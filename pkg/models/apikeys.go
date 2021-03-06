package models

import (
	"time"
)

// APIKey is an api key pair generated by the user to access the api without the frontend.
// swagger:model
type APIKey struct {
	// ID of the key. This is auto-generated.
	//
	// required: true
	ID int `json:"id"`
	// Name of the key
	//
	// required: true
	Name string `json:"name,omitempty"`
	// UserID is the ID of the user to which this key belongs.
	//
	// required: true
	UserID int `json:"user_id"`
	// APIKeyID is a generated id of the key.
	//
	// required: true
	APIKeyID string `json:"api_key_id"`
	// APIKeySecret is a generated secret, aka, the key.
	//
	// required: true
	APIKeySecret string `json:"api_key_secret"`
	// TTL defines how long this key can live in duration.
	//
	// required: true
	// example: 1h10m10s
	TTL string `json:"ttl"`
	// CreateAt defines when this key was created.
	//
	// required: true
	// example: time.Now()
	CreateAt time.Time `json:"create_at"`
}

// APIKeyAuthRequest contains a user email and their api key.
type APIKeyAuthRequest struct {
	Email        string `json:"email"`
	APIKeyID     string `json:"api_key_id"`
	APIKeySecret string `json:"api_key_secret"`
}
