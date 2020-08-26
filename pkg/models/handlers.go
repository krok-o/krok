package models

const (
	// All the different types of hooks.
	GITHUB = iota
	// GitLab based hooks
	GITLAB
	// Gitea based hooks
	GITEA
	// BitBucket based hooks
	BITBUCKET
)
