package models

type UserTokenProfile struct {
	UserID string
}

type UserAuthDetails struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
}
