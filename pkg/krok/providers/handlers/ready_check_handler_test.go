package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
)

func TestReadyCheckHandler_Ready(t *testing.T) {
	mch := &mocks.Ready{}
	mch.On("Ready", mock.Anything).Return(true)
	logger := zerolog.New(os.Stderr)
	checker := NewReadyCheckHandler(ReadyCheckHandlerDependencies{
		Logger:  logger,
		Checker: mch,
	})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := checker.Ready()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestReadyCheckHandler_NotReady(t *testing.T) {
	mch := &mocks.Ready{}
	mch.On("Ready", mock.Anything).Return(false)
	logger := zerolog.New(os.Stderr)
	checker := NewReadyCheckHandler(ReadyCheckHandlerDependencies{
		Logger:  logger,
		Checker: mch,
	})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := checker.Ready()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
