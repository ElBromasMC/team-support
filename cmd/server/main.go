package main

import (
	"alc/handler"
	middle "alc/middleware"
	"context"
	"log"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Database connection
	dbconfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Failed to parse config:", err)
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}
	defer dbpool.Close()

	// Initialize handler
	h := handler.Handler{
		DB: dbpool,
	}

	// Static files
	static(e)

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middle.Auth(dbpool))

	// Page routes
	e.GET("/", h.HandleIndexShow)
	e.GET("/ticket", h.HandleTicketShow)

	// Garantia routes
	e.GET("/garantia", h.HandleGarantiaShow)
	e.GET("/garantia/:slug", h.HandleGarantiaCategoryShow)
	e.GET("/garantia/:categorySlug/:itemSlug", h.HandleGarantiaItemShow)

	// Store routes
	e.GET("/store", h.HandleStoreShow)
	e.GET("/store/:slug", h.HandleStoreItemShow)

	// Cart routes
	e.GET("/cart", h.HandleCartShow)

	// User routes
	e.GET("/login", h.HandleLoginShow)
	e.GET("/signup", h.HandleSignupShow)
	e.POST("/login", h.HandleLogin)
	e.POST("/signup", h.HandleSignup)
	e.GET("/logout", h.HandleLogout)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(e.Start(":" + port))
}
