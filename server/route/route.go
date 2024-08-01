package route

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mdm-code/tq/toml"
	"github.com/mdm-code/tqweb/server/template"
)

// UseRootRoutes is a router grouping function.
func UseRootRoutes(e *echo.Echo) {
	e.GET("/", Index)
	e.Static("/", "static")
}

// UseInputRoutes groups together input processing handlers.
func UseInputRoutes(e *echo.Echo) {
	e.POST("/toml", CheckTOML)
}

// Index returns the landing page.
func Index(c echo.Context) error {
	component := template.Page()
	err := component.Render(c.Request().Context(), c.Response().Writer)
	return err
}

// CheckTOML verifies if the TOML data is syntactically correct.
func CheckTOML(c echo.Context) error {
	tomlAdapter := toml.NewAdapter(
		toml.NewGoTOML(
			toml.GoTOMLConf{},
		),
	)
	var tomlData any
	value := c.FormValue("tomlData")
	data := strings.NewReader(value)
	err := tomlAdapter.Unmarshal(data, &tomlData)
	var errMsg string
	var ot template.OutlineType
	if err != nil {
		errMsg = err.Error()
		ot = template.Error
	} else {
		ot = template.Correct
	}
	t := template.TomlInput(value, errMsg, ot)
	return t.Render(c.Request().Context(), c.Response().Writer)
}
