package models

// VCSToken represents a token for a platform.
type VCSToken struct {
	Token string `json:"token"`
	VCS   int    `json:"vcs"`
}
