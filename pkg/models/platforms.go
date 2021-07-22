package models

const (
	// All the different types of hooks.

	// GITHUB based hooks
	GITHUB = iota + 1
	// GITLAB based hooks
	GITLAB
	// GITEA based hooks
	GITEA
)

// Platform defines a platform like Github, Gitlab etc.
// swagger:model
type Platform struct {
	// ID of the platform. This is choosen.
	//
	// required: true
	ID int `json:"id"`
	// Name of the platform.
	//
	// required: true
	// example: github, gitlab, gitea
	Name string `json:"name"`
}

// SupportedPlatforms a map of supported platforms by Krok.
var SupportedPlatforms = map[int]Platform{
	GITHUB: {
		ID:   GITHUB,
		Name: "github",
	},
	GITLAB: {
		ID:   GITLAB,
		Name: "gitlab",
	},
	GITEA: {
		ID:   GITEA,
		Name: "gitea",
	},
}
