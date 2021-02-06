package providers

import (
	"github.com/labstack/echo/v4"
)

// RepositoryHandler defines the handler's capabilities.
// The handler is a front wrapper for database operations, but also provides
// additional abilities, i.e.: generate a unique url
type RepositoryHandler interface {
	CreateRepository() echo.HandlerFunc
	DeleteRepository() echo.HandlerFunc
	GetRepository() echo.HandlerFunc
	ListRepositories() echo.HandlerFunc
	UpdateRepository() echo.HandlerFunc
}

// CommandHandler defines the actions of commands.
type CommandHandler interface {
	DeleteCommand() echo.HandlerFunc
	ListCommands() echo.HandlerFunc
	GetCommand() echo.HandlerFunc
	UpdateCommand() echo.HandlerFunc

	// Relationship operations.

	// AddCommandRelForRepository adds an entry for this command id to the given repositoryID.
	AddCommandRelForRepository() echo.HandlerFunc
	// RemoveCommandRelForRepository remove a relation to a repository for a command.
	RemoveCommandRelForRepository() echo.HandlerFunc
}

// ApiKeysHandler provides functions which define operations on api key pairs.
type ApiKeysHandler interface {
	CreateApiKeyPair() echo.HandlerFunc
	DeleteApiKeyPair() echo.HandlerFunc
	ListApiKeyPairs() echo.HandlerFunc
	GetApiKeyPair() echo.HandlerFunc
}

// TokenHandler provides operations to get and validation JWT tokens.
type TokenHandler interface {
	TokenHandler() echo.HandlerFunc
}

// AuthHandler provides the handler functions for the authentication flow.
type AuthHandler interface {
	OAuthLogin() echo.HandlerFunc
	OAuthCallback() echo.HandlerFunc
	Refresh() echo.HandlerFunc
}
