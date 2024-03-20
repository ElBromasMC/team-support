package main

import (
	"alc/handler"
	middle "alc/middleware"
	"context"
	"log"
	"net/http"
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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	authMiddleware := middle.Auth(dbpool)

	// Static files
	static(e)

	// Images routes
	e.Static("/images", "images")

	// Page routes
	e.GET("/", h.HandleIndexShow, authMiddleware)
	e.GET("/ticket", h.HandleTicketShow, authMiddleware)

	// Garantia routes
	g1 := e.Group("/garantia")
	g1.Use(authMiddleware)
	g1.GET("", h.HandleGarantiaShow)
	g1.GET("/:slug", h.HandleGarantiaCategoryShow)
	g1.GET("/:categorySlug/:itemSlug", h.HandleGarantiaItemShow)

	// Store routes
	g2 := e.Group("/store")
	g2.Use(authMiddleware)
	g2.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/store/all")
	})
	g2.GET("/all", h.HandleStoreAllShow)
	g2.GET("/:slug", h.HandleStoreCategoryShow)
	g2.POST("/all", h.HandleStoreAllItemsShow)
	g2.POST("/:slug", h.HandleStoreCategoryItemsShow)
	g2.GET("/:categorySlug/:itemSlug", h.HandleStoreItemShow)

	// Cart routes
	e.GET("/cart", h.HandleCartShow, authMiddleware)

	// User routes
	e.GET("/login", h.HandleLoginShow)
	e.GET("/signup", h.HandleSignupShow)
	e.POST("/login", h.HandleLogin)
	e.POST("/signup", h.HandleSignup)
	e.GET("/logout", h.HandleLogout, authMiddleware)

	// Admin group
	g3 := e.Group("/admin")
	g3.Use(authMiddleware, middle.Admin)
	g3.GET("", h.HandleAdminShow)
	g3.GET("/garantia", h.HandleAdminGarantiaShow)
	g3.POST("/garantia", h.HandleNewGarantia)
	g3.PUT("/garantia", h.HandleUpdateGarantia)
	g3.GET("/store", h.HandleAdminStoreShow)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(e.Start(":" + port))
}
