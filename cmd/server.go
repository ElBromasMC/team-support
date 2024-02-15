package main

import (
	"alc/assets"
	"alc/handler"
	middle "alc/middleware"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Live reload
	http.Get("http://localhost:8020")

	e := echo.New()

	// Root level middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static files
	e.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(assets.Assets))))

	// Database connection
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer dbpool.Close()

	// Initialize handler
	h := handler.Handler{
		DB: dbpool,
	}

	// Auth middleware
	e.Use(middle.Auth(dbpool))

	// Routes
	e.GET("/", h.HandleIndexShow)
	e.GET("/store", h.HandleStoreShow)
	e.GET("/signup", h.HandleSignupShow)
	e.POST("/signup", h.HandleSignup)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
