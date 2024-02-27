//go:build !dev

package main

import (
	"alc/assets"
	"net/http"

	"github.com/labstack/echo/v4"
)

func static(e *echo.Echo) *echo.Route {
	return e.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(assets.Assets))))
}
