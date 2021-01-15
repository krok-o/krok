package models

const (
	// All the different types of hooks.

	// GITHUB based hooks
	GITHUB = iota + 1
	// GITLAB based hooks
	GITLAB
	// GITEA based hooks
	GITEA
	// BITBUCKET based hooks
	BITBUCKET
)

// VCS is a map for mapping a VCS string value to int.
var VCS = map[string]int{
	"GITHUB":    1,
	"GITLAB":    2,
	"GITEA":     3,
	"BITBUCKET": 4,
}
