package route

import (
	"github.com/labstack/echo/v4"
	"github.com/mdm-code/tqweb/server/component"
)

const (
	staticFsRoot = "assets"
)

// UseRootRoutes is a router grouping function.
func UseRootRoutes(e *echo.Echo) {
	e.Static("/assets", staticFsRoot)
	e.GET("/", Index)
}

// Index route for the tqweb.
func Index(c echo.Context) error {
	index := component.Index()
	err := index.Render(c.Request().Context(), c.Response().Writer)
	return err
}
