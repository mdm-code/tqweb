package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Server provides a dummy HTTP server.
func Server() *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		payload := map[string]string{
			"detail": "OK",
		}
		return c.JSON(http.StatusOK, payload)
	})
	return e
}
