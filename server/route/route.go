package route

import (
	"github.com/labstack/echo/v4"
)

// UseRootRoutes is a router grouping function.
func UseRootRoutes(e *echo.Echo) {
	e.Static("/", "static")
}
