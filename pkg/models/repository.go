package models

import "errors"

// Auth is authentication option for a repository.
// swagger:model
type Auth struct {
	// SSH private key.
	//
	// required: false
	SSH string `json:"ssh,omitempty"`
	// Username is the username required to access the platform for this repositroy.
	//
	// required: false
	Username string `json:"username,omitempty"`
	// Password is the password required to access the platform for this repositroy.
	//
	// required: false
	Password string `json:"password,omitempty"`
	// Hook secret to create a hook with on the respective platform.
	//
	// required: true
	Secret string `json:"secret"`
}

// GitLab contains gitLab specific settings.
// swagger:model
type GitLab struct {
	// ProjectID is an optional ID which defines a project in Gitlab.
	//
	// required: false
	ProjectID int `json:"project_id,omitempty"`
}

// GetProjectID returns the project id.
func (g *GitLab) GetProjectID() int {
	if g != nil {
		return g.ProjectID
	}
	return -1
}

// Repository is a repository which can be managed by Krok.
// swagger:model
type Repository struct {
	// Name of the repository.
	//
	// required: false
	Name string `json:"name"`
	// ID of the repository. Auto-generated.
	//
	// required: true
	ID int `json:"id"`
	// URL of the repository.
	//
	// required: true
	URL string `json:"url"`
	// VCS Defines which handler will be used. For values, see platforms.go.
	//
	// required: true
	VCS int `json:"vcs"`
	// GitLab specific settings.
	//
	// required: false
	GitLab *GitLab `json:"git_lab,omitempty"`
	// Auth contains authentication details for this repository.
	//
	// required: true
	Auth *Auth `json:"auth,omitempty"`
	// Commands contains all the commands which this repository is attached to.
	//
	// required: false
	Commands []*Command `json:"commands,omitempty"`
	// This field is not saved in the DB but generated every time the repository
	// details needs to be displayed.
	//
	// required: true
	UniqueURL string `json:"unique_url,omitempty"`
	// TODO: Think about storing this
	Events []string `json:"events,omitempty"`
}

// Validate validates this model.
func (r *Repository) Validate() (ok bool, field string, err error) {
	if r.Name == "" {
		return false, "Name", errors.New("name cannot be empty")
	}
	if r.Auth == nil {
		return false, "Auth", errors.New("auth information must be defined")
	}
	if r.Auth.Secret == "" {
		return false, "Auth.Secret", errors.New("secret for the webhook must be defined")
	}
	return true, "", nil
}
