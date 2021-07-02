package handlers

import (
	"net/http"

	"github.com/krok-o/krok/pkg/models"
	"github.com/labstack/echo/v4"
)

// SupportedPlatformList is the supporting struct which implements platform listing.
type SupportedPlatformList struct{}

// NewSupportedPlatformListHandler creates a new supported platform list handler.
func NewSupportedPlatformListHandler() *SupportedPlatformList {
	return &SupportedPlatformList{}
}

// ListSupportedPlatforms lists all platforms which Krok supports.
// swagger:operation GET /supported-platforms listSupportedPlatforms
// Lists all supported platforms.
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: 'the list of supported platform ids'
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Platform"
func (s *SupportedPlatformList) ListSupportedPlatforms() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.SupportedPlatforms)
	}
}
