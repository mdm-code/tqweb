package route

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mdm-code/tq"
	"github.com/mdm-code/tq/toml"
	"github.com/mdm-code/tqweb/server/component"
)

const (
	staticFsRoot     = "assets"
	staticPathPrefix = "/assets"
)

// RegisterAll registers all routes defined for the HTTP server.
func RegisterAll(e *echo.Echo) *echo.Echo {
	return ServeStatics(
		RegsiterRootRoutes(
			RegisterProcessRoutes(
				e,
			),
		),
		staticFsRoot,
		staticPathPrefix,
	)
}

// ServeStatics registers a new route to serve static files from the root
// directory at the provided path prefix.
func ServeStatics(e *echo.Echo, at, from string) *echo.Echo {
	e.Static(at, from)
	return e
}

// RegsiterRootRoutes groups root routes.
func RegsiterRootRoutes(e *echo.Echo) *echo.Echo {
	e.GET("/", Index)
	return e
}

// RegisterProcessRoutes groups data processing routes.
func RegisterProcessRoutes(e *echo.Echo) *echo.Echo {
	g := e.Group("/api/v1")
	g.POST("/inputData", ProcessInputData)
	g.POST("/tqQuery/validate", ValidateTqQuery)
	g.POST("/toml/validate", ValidateTOML)
	return e
}

// Index route for the tqweb.
func Index(c echo.Context) error {
	index := component.Index()
	err := index.Render(c.Request().Context(), c.Response().Writer)
	return err
}
