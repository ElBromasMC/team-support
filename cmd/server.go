package main

import (
	"alc/assets"
	"alc/handler"
	middle "alc/middleware"
	"context"
	"fmt"
	"net/http"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// Live reload
	http.Get("http://localhost:8020")

	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database connection
	dbconfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer dbpool.Close()

	// Initialize handler
	h := handler.Handler{
		DB: dbpool,
	}

	// Static files
	e.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(assets.Assets))))

	// Auth middleware
	e.Use(middle.Auth(dbpool))

	// Routes
	e.GET("/", h.HandleIndexShow)
	e.GET("/ticket", h.HandleTicketShow)
	e.GET("/garantia", h.HandleGarantiaShow)
	e.GET("/garantia/:slug", h.HandleGarantiaItemShow)

	e.GET("/store", h.HandleStoreShow)

	e.GET("/signup", h.HandleSignupShow)
	e.GET("/login", h.HandleLoginShow)
	e.POST("/signup", h.HandleSignup)
	e.POST("/login", h.HandleLogin)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
