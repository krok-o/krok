package models

// User is a user in the Krok system.
type User struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}
