package providers

import (
	"github.com/labstack/echo/v4"
)

// RepositoryHandler defines the handler's capabilities.
// The handler is a front wrapper for database operations, but also provides
// additional abilities, i.e.: generate a unique url
type RepositoryHandler interface {
	Create() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Get() echo.HandlerFunc
	List() echo.HandlerFunc
	Update() echo.HandlerFunc
}

// CommandHandler defines the actions of commands.
type CommandHandler interface {
	Delete() echo.HandlerFunc
	Get() echo.HandlerFunc
	List() echo.HandlerFunc
	Update() echo.HandlerFunc

	// Relationship operations.

	// AddCommandRelForRepository adds an entry for this command id to the given repositoryID.
	AddCommandRelForRepository() echo.HandlerFunc
	// RemoveCommandRelForRepository remove a relation to a repository for a command.
	RemoveCommandRelForRepository() echo.HandlerFunc

	// AddCommandRelForPlatform adds an entry for this command id to the given platform id.
	AddCommandRelForPlatform() echo.HandlerFunc
	// RemoveCommandRelForPlatform remove a relation to a platform for a command.
	RemoveCommandRelForPlatform() echo.HandlerFunc
}

// APIKeysHandler provides functions which define operations on api key pairs.
type APIKeysHandler interface {
	Create() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Get() echo.HandlerFunc
	List() echo.HandlerFunc
}

// TokenHandler provides operations to get and validation JWT tokens.
type TokenHandler interface {
	TokenHandler() echo.HandlerFunc
}

// VCSTokenHandler provides operations to manage tokens for the various platforms..
type VCSTokenHandler interface {
	Create() echo.HandlerFunc
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

// HookHandler represents what the Krok server is capable off.
type HookHandler interface {
	// HandleHooks handles all hooks incoming to Krok.
	HandleHooks() echo.HandlerFunc
}

// CommandSettingsHandler defines the actions of command settings.
type CommandSettingsHandler interface {
	Create() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Get() echo.HandlerFunc
	List() echo.HandlerFunc
	Update() echo.HandlerFunc
}

// SupportedPlatformListHandler lists all supported platforms.
type SupportedPlatformListHandler interface {
	ListSupportedPlatforms() echo.HandlerFunc
}
