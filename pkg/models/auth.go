package models

// UserAuthDetails represents the authenticated user details.
// swagger:model
type UserAuthDetails struct {
	// Email is the email of the registered user.
	//
	// required: true
	Email string
	// FirstName is the first name of the user.
	//
	// required: true
	FirstName string
	// LastName is the last name of the user.
	//
	// required: true
	LastName string
}

// TokenResponse contains the generated JWT token.
// swagger:response TokenResponse
type TokenResponse struct {
	// The generated token
	// in: body
	Token string
}
