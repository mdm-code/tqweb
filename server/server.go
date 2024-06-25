package server

import (
	"github.com/labstack/echo/v4"
	"github.com/mdm-code/tqweb/server/route"
)

// Server provides a dummy HTTP server.
func Server() *echo.Echo {
	e := echo.New()
	route.UseRootRoutes(e)
	return e
}
