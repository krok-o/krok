package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// UserAuthenticationConfig is the config for the UserAuthentication middleware.
type UserAuthenticationConfig struct {
	GlobalTokenKey string
	CookieName     string
}

// UserContext represents the user context.
type UserContext struct {
	UserID int
}

// GetUserContext gets the UserContext from the echo.Context.
// UserContext is created by the UserAuthentication middleware.
func GetUserContext(c echo.Context) (*UserContext, error) {
	user := c.Get("user")
	if user == nil {
		return nil, errors.New("user not found in context")
	}

	userContext, ok := user.(*UserContext)
	if !ok {
		return nil, errors.New("user not UserContext type")
	}

	return userContext, nil
}

// UserAuthentication catches the access_token and verifies it.
// We also extract the UserContext information from the JWT and set it in the echo.Context.
func UserAuthentication(cfg *UserAuthenticationConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := cfg.extractToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "failed to extract token")
			}

			var claims jwt.StandardClaims
			if _, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.GlobalTokenKey), nil
			}); err != nil {
				return c.JSON(http.StatusUnauthorized, "failed to verify token")
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "failed to get userID")
			}
			uc := &UserContext{UserID: userID}
			c.Set("user", uc)

			return next(c)
		}
	}
}

func (cfg *UserAuthenticationConfig) extractToken(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	token, err := c.Cookie(cfg.CookieName)
	if err != nil {
		return "", errors.New("failed to get token from header or cookie")
	}

	return token.Value, nil
}
