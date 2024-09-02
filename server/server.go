package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mdm-code/tqweb/server/route"
)

// Server provides a dummy HTTP server.
func Server() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e = route.RegisterAll(e)
	return e
}
