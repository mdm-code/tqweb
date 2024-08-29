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
	g.POST("/query/validate", ValidateTqQuery)
	g.POST("/toml/validate", ValidateTOML)
	return e
}

// Index route for the tqweb.
func Index(c echo.Context) error {
	index := component.Index()
	err := index.Render(c.Request().Context(), c.Response().Writer)
	return err
}

// ProcessInputData runs the tq query against the provided TOML data.
func ProcessInputData(c echo.Context) error {
	query := c.FormValue("tqQuery")
	tomlData := c.FormValue("tomlData")

	input := strings.NewReader(tomlData)
	var output bytes.Buffer

	// TODO: Config data should be provided through the form as well.
	tomlAdapter := toml.NewAdapter(toml.NewGoTOML(toml.GoTOMLConf{}))
	tq := tq.New(tomlAdapter)
	if err := tq.Run(input, &output, query); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "Unprocessable entity",
			Internal: err,
		}
	}
	// TODO: Otherwise try to render the output HTML with output.String().

	// TODO: Extend output TOML validation.
	_ = func(output string) bool {
		var data any
		r := strings.NewReader(output)
		if err := tomlAdapter.Unmarshal(r, &data); err != nil {
			return false
		}
		return true
	}(output.String())

	return nil
}

// ValidateTqQuery verifies if the provided tq query string is valid.
func ValidateTqQuery(c echo.Context) error {
	query := c.FormValue("tqQuery")
	tomlAdapter := toml.NewAdapter(toml.NewGoTOML(toml.GoTOMLConf{}))
	tq := tq.New(tomlAdapter)
	if err := tq.Validate(query); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "Unprocessable entity",
			Internal: err,
		}
	}
	// TODO: Otherwise try to render the output HTML.
	return nil
}

// ValidateTOML checks if the provided form input is a valid TOML document.
func ValidateTOML(c echo.Context) error {
	tomlData := c.FormValue("tomlData")
	tomlAdapter := toml.NewAdapter(toml.NewGoTOML(toml.GoTOMLConf{}))
	var data any
	reader := strings.NewReader(tomlData)
	if err := tomlAdapter.Unmarshal(reader, &data); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "Unprocessable entity",
			Internal: err,
		}
	}
	// TODO: Otherwise try to render the output HTML.
	return nil
}
