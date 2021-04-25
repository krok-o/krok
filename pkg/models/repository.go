package models

import "errors"

// Auth is authentication option for a repository.
type Auth struct {
	SSH      string `json:"ssh,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	// Hook secret to create a hook with on the respective platform.
	Secret string `json:"secret"`
}

// Repository is a repository which can be managed by Krok.
type Repository struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	URL  string `json:"url"`
	// Defines which handler will be used. For values, see platforms.go.
	VCS int `json:"vcs"`
	// ProjectID is an optional ID which defines a project in Gitlab.
	ProjectID *int `json:"project_id,omitempty"`
	// Auth an command are all dynamically generated.
	Auth     *Auth      `json:"auth,omitempty"`
	Commands []*Command `json:"commands,omitempty"`
	// This field is not saved in the DB but generated every time the repository
	// details needs to be displayed.
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
