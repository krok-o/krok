package models

// Auth is authentication option for a repository.
type Auth struct {
	SSH      string `json:"ssh,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Repository is a repository which can be managed by Krok.
type Repository struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	URL  string `json:"url"`
	// Defines which handler will be used. For values, see handlers.go.
	VCS int `json:"vcs"`
	// Auth an command are all dynamically generated.
	Auth     *Auth      `json:"auth,omitempty"`
	Commands []*Command `json:"commands,omitempty"`
	// This field is not saved in the DB but generated every time the repository
	// details needs to be displayed.
	UniqueURL string `json:"unique_url,omitempty"`
}
