package handlers

import (
	"net/http"
	"sort"

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
		// Create a slice and sort it so the response is predictable.
		platforms := []models.Platform{}
		for _, v := range models.SupportedPlatforms {
			platforms = append(platforms, v)
		}
		sort.Slice(platforms, func(i, j int) bool {
			return platforms[i].ID < platforms[j].ID
		})
		return c.JSON(http.StatusOK, platforms)
	}
}
