package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type JWTAuthConfig struct {
	GlobalTokenKey string
	CookieName     string
}

func JWTAuthentication(cfg *JWTAuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := cfg.extractToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "failed to extract token")
			}

			var refreshClaims jwt.StandardClaims
			if _, err := jwt.ParseWithClaims(token, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.GlobalTokenKey), nil
			}); err != nil {
				return c.JSON(http.StatusUnauthorized, "failed to verify token")
			}

			return next(c)
		}
	}
}

func (cfg *JWTAuthConfig) extractToken(c echo.Context) (string, error) {
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
