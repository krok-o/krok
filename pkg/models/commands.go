package models

// Command is a command which can be executed by Krok.
type Command struct {
	Name         string        `json:"name"`
	ID           int           `json:"id"`
	Schedule     string        `json:"schedule,omitempty"`
	Repositories []*Repository `json:"repositories,omitempty"`
	Filename     string        `json:"filename"`
	Location     string        `json:"location"`
	Hash         string        `json:"hash"`
	Enabled      bool          `json:"enabled"`
}

// CommandSetting defines the settings a command can have.
type CommandSetting struct {
	ID        int    `json:"id"`
	CommandID int    `json:"command_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	InVault   bool   `json:"in_vault"`
}
