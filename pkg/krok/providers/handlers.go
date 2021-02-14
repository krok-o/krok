package providers

import (
	"github.com/labstack/echo/v4"
)

// CRUDHandler defines basic crud operations for a resource.
type CRUDHandler interface {
	Create() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Get() echo.HandlerFunc
	List() echo.HandlerFunc
	Update() echo.HandlerFunc
}

// RepositoryHandler defines the handler's capabilities.
// The handler is a front wrapper for database operations, but also provides
// additional abilities, i.e.: generate a unique url
type RepositoryHandler interface {
	CRUDHandler
}

// CommandHandler defines the actions of commands.
type CommandHandler interface {
	CRUDHandler
	// Relationship operations.

	// AddCommandRelForRepository adds an entry for this command id to the given repositoryID.
	AddCommandRelForRepository() echo.HandlerFunc
	// RemoveCommandRelForRepository remove a relation to a repository for a command.
	RemoveCommandRelForRepository() echo.HandlerFunc
}

// ApiKeysHandler provides functions which define operations on api key pairs.
type ApiKeysHandler interface {
	CRUDHandler
}

// TokenHandler provides operations to get and validation JWT tokens.
type TokenHandler interface {
	TokenHandler() echo.HandlerFunc
}

// VCSTokenHandler provides operations to manage tokens for the various platforms..
type VCSTokenHandler interface {
	CRUDHandler
}

// UserMiddleware provides UserMiddleware authentication capabilities.
type UserMiddleware interface {
	JWT() echo.MiddlewareFunc
}

// UserTokenHandler provides operations for user personal access tokens.
type UserTokenHandler interface {
	Generate() echo.HandlerFunc
}

// AuthHandler provides the handler functions for the authentication flow.
type AuthHandler interface {
	OAuthLogin() echo.HandlerFunc
	OAuthCallback() echo.HandlerFunc
	Refresh() echo.HandlerFunc
}
