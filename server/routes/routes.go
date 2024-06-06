package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// UseRootRoutes is a router grouping function.
func UseRootRoutes(e *echo.Echo) {
	e.GET("/", Index)
	e.Static("/", "static")
}

// Index returns the landing page.
func Index(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
  <script src="/js/htmx.js"></script>
</head>
<body>
  <title>Hello, world!</title>
</body>
`
	return c.HTML(http.StatusOK, html)
}
