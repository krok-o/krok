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
