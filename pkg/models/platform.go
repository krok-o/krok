package models

// Platform defines a platform like Github, Gitlab etc.
type Platform struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
