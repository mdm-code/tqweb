package route

import (
	"net/http"

	_ "github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// UseRootRoutes is a router grouping function.
func UseRootRoutes(e *echo.Echo) {
	e.Static("/", "static")
	e.GET("/", Index)
}

// Index route for the application.
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
