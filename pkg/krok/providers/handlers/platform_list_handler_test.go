package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSupportedPlatformListHandler(t *testing.T) {
	handler := NewSupportedPlatformListHandler()
	expectedSupportPlatformList := `{"list":[{"id":1,"name":"github"},{"id":2,"name":"gitlab"},{"id":3,"name":"gitea"}]}
`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/supported-platforms", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := handler.ListSupportedPlatforms()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedSupportPlatformList, rec.Body.String())
}
