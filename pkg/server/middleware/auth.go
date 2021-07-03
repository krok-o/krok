package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const (
	userContextKey = "user"
	// The RFC defines this value as case sensitive.
	// We acknowledge the RFC and will reject requests which use this
	// as a case insensitive value.
	// RFC: https://tools.ietf.org/html/rfc6750 section 1.1. Notational Conventions
	bearerHeader = "Bearer "
)

// UserMiddlewareConfig represents the UserMiddleware config.
type UserMiddlewareConfig struct {
	GlobalTokenKey string
	CookieName     string
}

// UserMiddlewareDeps represents the UserMiddleware dependencies.
type UserMiddlewareDeps struct {
	Logger    zerolog.Logger
	UserStore providers.UserStorer
}

// UserMiddleware represents our user middleware.
type UserMiddleware struct {
	UserMiddlewareConfig
	UserMiddlewareDeps
}

// NewUserMiddleware creates a new UserMiddleware.
func NewUserMiddleware(cfg UserMiddlewareConfig, deps UserMiddlewareDeps) *UserMiddleware {
	return &UserMiddleware{UserMiddlewareConfig: cfg, UserMiddlewareDeps: deps}
}

// UserContext represents the user context.
type UserContext struct {
	UserID int
}

// GetUserContext gets the UserContext from the echo.Context.
// UserContext is created by the UserAuthentication middleware.
func GetUserContext(c echo.Context) (*UserContext, error) {
	user := c.Get(userContextKey)
	if user == nil {
		return nil, errors.New("user not found in context")
	}

	userContext, ok := user.(*UserContext)
	if !ok {
		return nil, errors.New("user not UserContext type")
	}

	return userContext, nil
}

// JWT catches the access_token and verifies it.
// We also extract the UserContext information from the JWT and set it in the echo.Context.
func (um *UserMiddleware) JWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := um.extractToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "failed to extract token")
			}

			var claims jwt.StandardClaims
			if _, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(um.GlobalTokenKey), nil
			}); err != nil {
				um.Logger.Warn().Err(err).Msg("jwt token authentication failed")
				return c.JSON(http.StatusUnauthorized, "Token authentication failed.")
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				um.Logger.Warn().Err(err).Msg("failed to parse subject to userID")
				return c.JSON(http.StatusInternalServerError, "Unexpected error.")
			}
			um.setUser(c, userID)

			return next(c)
		}
	}
}

func (um *UserMiddleware) setUser(c echo.Context, userID int) {
	uc := &UserContext{UserID: userID}
	c.Set(userContextKey, uc)
}

func (um *UserMiddleware) extractToken(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, bearerHeader) {
		return strings.TrimPrefix(authHeader, bearerHeader), nil
	}

	token, err := c.Cookie(um.CookieName)
	if err != nil {
		return "", errors.New("failed to get token from header or cookie")
	}

	return token.Value, nil
}
