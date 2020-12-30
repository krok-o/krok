package providers

import (
	"github.com/labstack/echo/v4"
)

// RepositoriesHandler defines the handler's capabilities.
// The handler is a front wrapper for database operations, but also provides
// additional abilities, i.e.: generate a unique url
type RepositoriesHandler interface {
	Create() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Get() echo.HandlerFunc
	List() echo.HandlerFunc
}
