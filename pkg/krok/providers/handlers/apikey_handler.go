package handlers

import (
	"github.com/labstack/echo/v4"

	"github.com/krok-o/krok/pkg/krok/providers"
)

// ApiKeysHandlerDependencies defines the dependencies for the api keys handler provider.
type ApiKeysHandlerDependencies struct {
	Dependencies
	APIKeysStore  providers.APIKeysStorer
	TokenProvider *TokenProvider
}

// ApiKeysHandler is a handler taking care of api keys related api calls.
type ApiKeysHandler struct {
	Config
	ApiKeysHandlerDependencies
}

var _ providers.ApiKeysHandler = &ApiKeysHandler{}

// NewApiKeysHandler creates a new api key pair handler.
func NewApiKeysHandler(cfg Config, deps ApiKeysHandlerDependencies) (*ApiKeysHandler, error) {
	return &ApiKeysHandler{
		Config:                     cfg,
		ApiKeysHandlerDependencies: deps,
	}, nil
}

func (a ApiKeysHandler) CreateApiKeyPair() echo.HandlerFunc {
	panic("implement me")
}

func (a ApiKeysHandler) DeleteApiKeyPair() echo.HandlerFunc {
	panic("implement me")
}

func (a ApiKeysHandler) ListApiKeyPairs() echo.HandlerFunc {
	panic("implement me")
}

func (a ApiKeysHandler) GetApiKeyPair() echo.HandlerFunc {
	panic("implement me")
}
