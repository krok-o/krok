package models

// Command is a command which can be executed by Krok.
type Command struct {
	Name         string        `json:"name"`
	ID           string        `json:"id"`
	Schedule     string        `json:"schedule,omitempty"`
	Repositories []*Repository `json:"repositories,omitempty"`
	Filename     string        `json:"filename"`
	Location     string        `json:"location"`
	Hash         string        `json:"hash"`
	Enabled      bool          `json:"enabled"`
}
