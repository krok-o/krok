package models

// Auth is authentication option for a repository.
type Auth struct {
	SSH      string `json:"ssh,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Repository is a repository which can be managed by Krok.
type Repository struct {
	Name     string     `json:"name"`
	ID       int        `json:"id"`
	URL      string     `json:"url"`
	Auth     *Auth      `json:"auth,omitempty"`
	Commands []*Command `json:"commands,omitempty"`
}
