package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mdm-code/tqweb/server/templates"
)

// UseRootRoutes is a router grouping function.
func UseRootRoutes(e *echo.Echo) {
	e.GET("/", Index)
	e.GET("/data", Data)
	e.Static("/", "static")
}

// Index returns the landing page.
func Index(c echo.Context) error {
	component := templates.Index()
	err := component.Render(c.Request().Context(), c.Response().Writer)
	return err
}

// Data returns dummy data.
func Data(c echo.Context) error {
	p := templates.Person{Name: "Mike", Age: 22}
	compoent := templates.Data(p)
	err := compoent.Render(c.Request().Context(), c.Response().Writer)
	return err
}
