package providers

import (
	"github.com/labstack/echo/v4"
)

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
