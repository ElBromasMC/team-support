package main

import (
	"alc/handler"
	middle "alc/middleware"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Root level middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Database connection
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Static files
	e.Static("/static", "static")

	// Initialize handler
	h := handler.Handler{
		DB: db,
	}

	// Auth middleware
	e.Use(middle.Auth(db))

	// Routes
	e.GET("/", h.HandleIndexShow)
	e.GET("/store", h.HandleStoreShow)
	e.GET("/signup", h.HandleSignupShow)
	e.POST("/signup", h.HandleSignup)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
