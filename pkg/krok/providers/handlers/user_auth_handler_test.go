package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
)

func TestUserAuthHandler_Login(t *testing.T) {
	e := echo.New()

	t.Run("missing redirect_url returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := NewUserAuthHandler(UserAuthHandlerDeps{})
		err := handler.OAuthLogin()(c)
		assert.NoError(t, err)
		assert.Equal(t, "error invalid redirect_url", rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("generate state error returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?redirect_url=https://test.com", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuth := &mocks.OAuthAuthenticator{}
		mockOAuth.On("GenerateState", "https://test.com").Return("", errors.New("error"))

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
			OAuth:  mockOAuth,
		})
		err := handler.OAuthLogin()(c)
		mockOAuth.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "error generating state", rec.Body.String())
	})

	t.Run("successful call", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?redirect_url=https://test.com", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuth := &mocks.OAuthAuthenticator{}
		mockOAuth.On("GenerateState", "https://test.com").Return("fake_state", nil)
		mockOAuth.On("GetAuthCodeURL", "fake_state").Return("https://test.com/auth")

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
			OAuth:  mockOAuth,
		})
		err := handler.OAuthLogin()(c)
		mockOAuth.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "https://test.com/auth", rec.Result().Header.Get("Location"))
	})
}

func TestUserAuthHandler_Callback(t *testing.T) {
	e := echo.New()

	t.Run("missing state query param returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
		})
		err := handler.OAuthCallback()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "error invalid state", rec.Body.String())
	})

	t.Run("missing code query param returns 400", func(t *testing.T) {
		qp := url.Values{}
		qp.Set("state", "fake_state")
		req := httptest.NewRequest(http.MethodGet, "/?"+qp.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
		})
		err := handler.OAuthCallback()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "error invalid code", rec.Body.String())
	})

	t.Run("error verifying state returns 401", func(t *testing.T) {
		qp := url.Values{}
		qp.Set("state", "fake_state")
		qp.Set("code", "fake_code")
		req := httptest.NewRequest(http.MethodGet, "/?"+qp.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuth := &mocks.OAuthAuthenticator{}
		mockOAuth.On("VerifyState", "fake_state").Return("", errors.New("err"))

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
			OAuth:  mockOAuth,
		})
		err := handler.OAuthCallback()(c)
		mockOAuth.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "error verifying state", rec.Body.String())
	})

	t.Run("error during token exchange returns 401", func(t *testing.T) {
		qp := url.Values{}
		qp.Set("state", "fake_state")
		qp.Set("code", "fake_code")
		req := httptest.NewRequest(http.MethodGet, "/?"+qp.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuth := &mocks.OAuthAuthenticator{}
		mockOAuth.On("VerifyState", "fake_state").Return("https://test.com", nil)
		mockOAuth.On("Exchange", mock.Anything, "fake_code").Return(nil, errors.New("err"))

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
			OAuth:  mockOAuth,
		})
		err := handler.OAuthCallback()(c)
		mockOAuth.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "error during token exchange", rec.Body.String())
	})

	t.Run("error during token exchange returns 401", func(t *testing.T) {
		qp := url.Values{}
		qp.Set("state", "fake_state")
		qp.Set("code", "fake_code")
		req := httptest.NewRequest(http.MethodGet, "/?"+qp.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuth := &mocks.OAuthAuthenticator{}
		mockOAuth.On("VerifyState", "fake_state").Return("https://test.com", nil)
		mockOAuth.On("Exchange", mock.Anything, "fake_code").Return(nil, errors.New("err"))

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
			OAuth:  mockOAuth,
		})
		err := handler.OAuthCallback()(c)
		mockOAuth.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "error during token exchange", rec.Body.String())
	})

	t.Run("success redirects to url from state", func(t *testing.T) {
		qp := url.Values{}
		qp.Set("state", "fake_state")
		qp.Set("code", "fake_code")
		req := httptest.NewRequest(http.MethodGet, "/?"+qp.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuth := &mocks.OAuthAuthenticator{}
		mockOAuth.On("VerifyState", "fake_state").Return("https://test.com", nil)
		mockOAuth.On("Exchange", mock.Anything, "fake_code").Return(&oauth2.Token{
			AccessToken:  "aaa",
			RefreshToken: "rrr",
		}, nil)

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
			OAuth:  mockOAuth,
		})
		err := handler.OAuthCallback()(c)
		mockOAuth.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusPermanentRedirect, rec.Code)
		assert.Equal(t, "https://test.com", rec.Result().Header.Get("Location"))

		cookies := rec.Result().Cookies()
		assert.Equal(t, "_a_token_", cookies[0].Name)
		assert.Equal(t, "/", cookies[0].Path)
		assert.Equal(t, "aaa", cookies[0].Value)
		assert.Equal(t, http.SameSiteStrictMode, cookies[0].SameSite)
		assert.True(t, cookies[0].HttpOnly)
		assert.Equal(t, "_r_token_", cookies[1].Name)
		assert.Equal(t, "/", cookies[1].Path)
		assert.Equal(t, "rrr", cookies[1].Value)
		assert.Equal(t, http.SameSiteStrictMode, cookies[1].SameSite)
		assert.True(t, cookies[1].HttpOnly)
	})
}

func TestUserAuthHandler_Refresh(t *testing.T) {
	e := echo.New()

	t.Run("missing refresh token cookie returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger: zerolog.New(os.Stderr),
		})
		err := handler.Refresh()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "error getting refresh token", rec.Body.String())
	})

	t.Run("error refreshing token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.AddCookie(&http.Cookie{Name: "_r_token_", Value: "fake_token"})

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTokenIssuer := &mocks.TokenIssuer{}
		mockTokenIssuer.On("Refresh", mock.Anything, "fake_token").Return(nil, errors.New("err"))

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger:      zerolog.New(os.Stderr),
			TokenIssuer: mockTokenIssuer,
		})
		err := handler.Refresh()(c)
		mockTokenIssuer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "error refreshing token", rec.Body.String())
	})

	t.Run("refresh token success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.AddCookie(&http.Cookie{Name: "_r_token_", Value: "fake_token"})

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTokenIssuer := &mocks.TokenIssuer{}
		mockTokenIssuer.On("Refresh", mock.Anything, "fake_token").Return(&oauth2.Token{
			AccessToken:  "fake_access_token",
			RefreshToken: "fake_refresh_token",
		}, nil)

		handler := NewUserAuthHandler(UserAuthHandlerDeps{
			Logger:      zerolog.New(os.Stderr),
			TokenIssuer: mockTokenIssuer,
		})
		err := handler.Refresh()(c)
		mockTokenIssuer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "", rec.Body.String())

		cookies := rec.Result().Cookies()
		assert.Equal(t, "_a_token_", cookies[0].Name)
		assert.Equal(t, "/", cookies[0].Path)
		assert.Equal(t, "fake_access_token", cookies[0].Value)
		assert.Equal(t, http.SameSiteStrictMode, cookies[0].SameSite)
		assert.True(t, cookies[0].HttpOnly)
		assert.Equal(t, "_r_token_", cookies[1].Name)
		assert.Equal(t, "/", cookies[1].Path)
		assert.Equal(t, "fake_refresh_token", cookies[1].Value)
		assert.Equal(t, http.SameSiteStrictMode, cookies[1].SameSite)
		assert.True(t, cookies[1].HttpOnly)
	})
}
