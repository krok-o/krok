package models

// ListOptions provides options for List operations with additional filters.
type ListOptions struct {
	Name string
	// List all repositories for Git, Gitea...
	VCS int
	// Add paging later on.
}
