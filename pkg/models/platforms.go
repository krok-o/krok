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
type Platform struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SupportedPlatforms a list of supported platforms by Krok.
var SupportedPlatforms = []Platform{
	{
		ID:   GITHUB,
		Name: "github",
	},
	{
		ID:   GITLAB,
		Name: "gitlab",
	},
	{
		ID:   GITEA,
		Name: "gitea",
	},
}
