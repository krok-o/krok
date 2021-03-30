package handlers

import (
	"net/http"

	"github.com/krok-o/krok/pkg/models"
	"github.com/labstack/echo/v4"
)

type SupportedPlatformList struct{}

// NewSupportedPlatformListHandler creates a new supported platform list handler.
func NewSupportedPlatformListHandler() *SupportedPlatformList {
	return &SupportedPlatformList{}
}

type Message struct {
	List []models.Platform `json:"list"`
}

func (s *SupportedPlatformList) ListSupportedPlatforms() echo.HandlerFunc {
	m := Message{
		List: models.SupportedPlatforms,
	}
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, m)
	}
}
