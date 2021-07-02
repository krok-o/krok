package models

// VCSToken represents a token for a platform.
// swagger:model
type VCSToken struct {
	// Token is the actual token.
	//
	// required: true
	Token string `json:"token"`
	// VCS is the ID of the platform to which this token belongs to.
	//
	// required: true
	VCS int `json:"vcs"`
}
