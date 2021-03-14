package models

// ListOptions provides options for List operations with additional filters.
type ListOptions struct {
	Name string `json:"name,omitempty"`
	// List all repositories for Git, Gitea...
	VCS int `json:"vcs,omitempty"`
	// Current Page
	Page int `json:"page,omitempty"`
	// Items per Page
	PageSize int `json:"page_size,omitempty"`
}
