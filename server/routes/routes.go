package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mdm-code/tqweb/server/template"
)

// UseRootRoutes is a router grouping function.
func UseRootRoutes(e *echo.Echo) {
	e.GET("/", Index)
	e.Static("/", "static")
}

// Index returns the landing page.
func Index(c echo.Context) error {
	component := template.Page()
	err := component.Render(c.Request().Context(), c.Response().Writer)
	return err
}
