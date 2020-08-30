package models

// Auth is authentication option for a repositroy.
type Auth struct {
	SSH      string `json:"ssh"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Repository is a repository which can be managed by Krok.
type Repository struct {
	Name     string    `json:"name"`
	ID       string    `json:"id"`
	URL      string    `json:"url"`
	Auth     Auth      `json:"auth"`
	Commands []Command `json:"commands"`
}
