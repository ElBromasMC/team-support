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
	e.Use(middle.Auth(dbpool))

	// Routes
	e.GET("/", h.HandleIndexShow)
	e.GET("/ticket", h.HandleTicketShow)

	e.GET("/garantia", h.HandleGarantiaShow)
	e.GET("/garantia/:slug", h.HandleGarantiaItemShow)

	e.GET("/store", h.HandleStoreShow)
	e.GET("/store/:slug", h.HandleStoreItemShow)

	e.GET("/cart", h.HandleCartShow)

	e.GET("/signup", h.HandleSignupShow)
	e.GET("/login", h.HandleLoginShow)
	e.POST("/signup", h.HandleSignup)
	e.POST("/login", h.HandleLogin)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(e.Start(":" + port))
}
