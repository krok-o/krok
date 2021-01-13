package handlers

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	keyTTL = 7 * 24 * time.Hour
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

// CreateApiKeyPair creates an api key pair for a given user.
func (a *ApiKeysHandler) CreateApiKeyPair() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := a.TokenProvider.GetToken(c)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
		id := c.Param("uid")
		if id == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		uid, err := strconv.Atoi(id)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		name := c.Param("name")
		if name == "" {
			name = "My Api Key"
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()

		// generate the key secret
		// this will be displayed once, then never shown again, ever.
		secret, err := a.generateUniqueKey()
		if err != nil {
			apiError := kerr.APIError("failed to generate unique api key id", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		// generate the key id
		// this will be displayed once, then never shown again, ever.
		keyID, err := a.generateKeyID()
		if err != nil {
			apiError := kerr.APIError("failed to generate unique api key id", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		encrypted, err := a.ApiKeyAuth.Encrypt(ctx, []byte(secret))
		if err != nil {
			apiError := kerr.APIError("failed to encrypt key", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		key := &models.APIKey{
			Name:         name,
			UserID:       uid,
			APIKeyID:     keyID,
			APIKeySecret: encrypted,
			TTL:          time.Now().Add(keyTTL),
		}

		generatedKey, err := a.APIKeysStore.Create(ctx, key)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("Failed to generate a key.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to generate key", http.StatusInternalServerError, err))
		}
		// We will display the ID and the secret unencrypted so the user can save it.
		key.ID = generatedKey.ID
		key.APIKeySecret = []byte(secret)

		return c.JSON(http.StatusOK, key)
	}
}

// DeleteApiKeyPair deletes a set of api keys for a given user with a given id.
func (a *ApiKeysHandler) DeleteApiKeyPair() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: make sure other users, can't brute force delete other user's keys.
		_, err := a.TokenProvider.GetToken(c)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
		uid := c.Param("uid")
		if uid == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		un, err := strconv.Atoi(uid)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		kid := c.Param("keyid")
		if kid == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		kn, err := strconv.Atoi(kid)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		if err := a.APIKeysStore.Delete(ctx, kn, un); err != nil {
			a.Logger.Debug().Err(err).Msg("ApiKey Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete api key", http.StatusBadRequest, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// ListApiKeyPairs lists all api keys for a given user.
func (a *ApiKeysHandler) ListApiKeyPairs() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: Consider getting the user from the token. But that would require a call to the cache or the DB.
		_, err := a.TokenProvider.GetToken(c)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
		uid := c.Param("uid")
		if uid == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		un, err := strconv.Atoi(uid)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		list, err := a.APIKeysStore.List(ctx, un)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("ApiKeys List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list api keys", http.StatusBadRequest, err))
		}
		return c.JSON(http.StatusOK, list)
	}
}

// GetApiKeyPair returns a given api key.
func (a *ApiKeysHandler) GetApiKeyPair() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := a.TokenProvider.GetToken(c)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
		uid := c.Param("uid")
		if uid == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		un, err := strconv.Atoi(uid)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		kid := c.Param("keyid")
		if kid == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		kn, err := strconv.Atoi(kid)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		key, err := a.APIKeysStore.Get(ctx, kn, un)
		if err != nil {
			apiError := kerr.APIError("failed to get api key", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		return c.JSON(http.StatusOK, key)
	}
}

// Generate a unique api key for a user.
func (a *ApiKeysHandler) generateUniqueKey() (string, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}

	return u.String(), nil
}

// Generate a unique api key for a user.
func (a *ApiKeysHandler) generateKeyID() (string, error) {
	u, err := a.generateUniqueKey()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(u))), nil
}