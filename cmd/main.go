package main

import (
	"os"

	"fmt"

	"github.com/labstack/echo/v4"

	"alc/handler/pages"
)

func main() {
	e := echo.New()

	// Static files
	e.Static("/static", "static")

	// Routes
	indexHandler := pages.IndexHandler{}
	e.GET("/", indexHandler.HandleIndexShow)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
