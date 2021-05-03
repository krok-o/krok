package handlers

import (
	"errors"
	"fmt"
	"net/http"

	kerr "github.com/krok-o/krok/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// VaultHandlerDependencies defines the dependencies for the vault settings handler provider.
type VaultHandlerDependencies struct {
	Logger zerolog.Logger
	Vault  providers.Vault
}

// VaultHandler is a handler taking care of vault related api calls.
type VaultHandler struct {
	VaultHandlerDependencies
}

var _ providers.VaultHandler = &VaultHandler{}

// NewVaultHandler creates a new vault settings handler.
func NewVaultHandler(deps VaultHandlerDependencies) *VaultHandler {
	return &VaultHandler{
		VaultHandlerDependencies: deps,
	}
}

// GetSecret will return all information including the secret.
func (v *VaultHandler) GetSecret() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.JSON(http.StatusBadRequest, kerr.APIError("name parameter missing", http.StatusBadRequest, errors.New("parameter name missing")))
		}
		// open the vault
		if err := v.Vault.LoadSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to open vault", http.StatusInternalServerError, err))
		}
		value, err := v.Vault.GetSecret(name)
		if err != nil {
			return c.JSON(http.StatusNotFound, kerr.APIError("secret not found", http.StatusNotFound, err))
		}
		return c.JSON(http.StatusOK, &models.VaultSetting{
			Key:   name,
			Value: string(value),
		})
	}
}

// ListSecrets will return all settings, but not the values.
func (v *VaultHandler) ListSecrets() echo.HandlerFunc {
	return func(c echo.Context) error {
		// open the vault
		if err := v.Vault.LoadSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to open vault", http.StatusInternalServerError, err))
		}
		value := v.Vault.ListSecrets()
		return c.JSON(http.StatusOK, value)
	}
}

// DeleteSecret deletes secrets.
func (v *VaultHandler) DeleteSecret() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.JSON(http.StatusBadRequest, kerr.APIError("name parameter missing", http.StatusBadRequest, errors.New("parameter name missing")))
		}
		// open the vault
		if err := v.Vault.LoadSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to open vault", http.StatusInternalServerError, err))
		}
		// we get first in order to return something to the user.
		if _, err := v.Vault.GetSecret(name); err != nil {
			return c.JSON(http.StatusNotFound, kerr.APIError("secret not found", http.StatusNotFound, err))
		}

		v.Vault.DeleteSecret(name)
		if err := v.Vault.SaveSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to save the vault after delete", http.StatusInternalServerError, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// UpdateSecret will update a given secret.
func (v *VaultHandler) UpdateSecret() echo.HandlerFunc {
	return func(c echo.Context) error {
		var update *models.VaultSetting
		if err := c.Bind(&update); err != nil {
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind vault settings", http.StatusBadRequest, err))
		}
		// open the vault
		if err := v.Vault.LoadSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to open vault", http.StatusInternalServerError, err))
		}
		// we get first in order to return something to the user.
		va, err := v.Vault.GetSecret(update.Key)
		if err != nil {
			return c.JSON(http.StatusNotFound, kerr.APIError("secret not found", http.StatusNotFound, err))
		}
		fmt.Println("value: ", string(va))

		v.Vault.AddSecret(update.Key, []byte(update.Value))
		if err := v.Vault.SaveSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to save the vault after update", http.StatusInternalServerError, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// CreateSecret will create a new secret.
func (v *VaultHandler) CreateSecret() echo.HandlerFunc {
	return func(c echo.Context) error {
		var vaultSetting *models.VaultSetting
		if err := c.Bind(&vaultSetting); err != nil {
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind vault settings", http.StatusBadRequest, err))
		}
		// open the vault
		if err := v.Vault.LoadSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to open vault", http.StatusInternalServerError, err))
		}

		v.Vault.AddSecret(vaultSetting.Key, []byte(vaultSetting.Value))
		if err := v.Vault.SaveSecrets(); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to save the vault after creation", http.StatusInternalServerError, err))
		}
		return c.NoContent(http.StatusCreated)
	}
}
