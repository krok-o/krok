package handlers

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetParamAsInt returns a number for an echo context parameter
// and checks its existence.
func GetParamAsInt(name string, c echo.Context) (int, error) {
	param := c.Param(name)
	if param == "" {
		return 0, errors.New("parameter not found")
	}

	n, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return n, nil
}
