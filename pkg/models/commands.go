package models

// Command is a command which can be executed by Krok.
// swagger:model
type Command struct {
	// Name of the command.
	//
	// required: true
	Name string `json:"name"`
	// ID of the command. Generated.
	//
	// required: true
	ID int `json:"id"`
	// Schedule of the command.
	//
	// required: false
	// example: 0 * * * * // follows cron job syntax.
	Schedule string `json:"schedule,omitempty"`
	// Repositories that this command can execute on.
	//
	// required: false
	Repositories []*Repository `json:"repositories,omitempty"`
	// Image defines the image name and tag of the command
	// Note: At the moment, only docker is supported. Later, runc, containerd...
	//
	// required: true
	// example: krok-hook/slack-notification:v0.0.1
	Image string `json:"image"`
	// Enabled defines if this command can be executed or not.
	//
	// required: false
	// example: false
	Enabled bool `json:"enabled"`
	// Platforms holds all the platforms which this command supports.
	// Calculated, not saved.
	//
	// required: false
	Platforms []Platform `json:"providers,omitempty"`
}

// CommandSetting defines the settings a command can have.
// swagger:model
type CommandSetting struct {
	// ID is a generated ID.
	//
	// required: true
	ID int `json:"id"`
	// CommandID is the ID of the command to which these settings belong to.
	//
	// required: true
	CommandID int `json:"command_id"`
	// Key is the name of the setting.
	//
	// required: true
	Key string `json:"key"`
	// Value is the value of the setting.
	//
	// required: true
	Value string `json:"value"`
	// InVault defines if this is sensitive information and should be stored securely.
	//
	// required: false
	InVault bool `json:"in_vault"`
}
