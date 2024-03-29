package main

import (
	"alc/handler/admin"
	"alc/handler/public"
	middle "alc/middleware"
	"alc/service"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	if os.Getenv("ENV") == "development" {
		e.Debug = true
	}

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

	// Initialize services
	ps := service.Public{
		DB: dbpool,
	}
	as := service.Admin{
		Public: ps,
	}

	// Initialize handlers
	ph := public.Handler{
		PublicService: ps,
	}

	ah := admin.Handler{
		AdminService: as,
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		log.Fatalln("Missing SESSION_KEY env variable")
	}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(key))))
	authMiddleware := middle.Auth(dbpool)
	cartMiddleware := middle.Cart(ps)

	// Static files
	static(e)

	// Images routes
	e.Static("/images", "images")

	// Page routes
	e.GET("/", ph.HandleIndexShow, authMiddleware, cartMiddleware)
	e.GET("/ticket", ph.HandleTicketShow, authMiddleware, cartMiddleware)

	// Garantia routes
	g1 := e.Group("/garantia")
	g1.Use(authMiddleware, cartMiddleware)
	g1.GET("", ph.HandleGarantiaShow)
	g1.GET("/:slug", ph.HandleGarantiaCategoryShow)
	g1.GET("/:categorySlug/:itemSlug", ph.HandleGarantiaItemShow)

	// Store routes
	g2 := e.Group("/store")
	g2.Use(authMiddleware, cartMiddleware)
	g2.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/store/all")
	})
	g2.GET("/all", ph.HandleStoreAllShow)
	g2.GET("/:slug", ph.HandleStoreCategoryShow)
	g2.POST("/all", ph.HandleStoreAllItemsShow)
	g2.POST("/:slug", ph.HandleStoreCategoryItemsShow)
	g2.GET("/:categorySlug/:itemSlug", ph.HandleStoreItemShow)

	// User routes
	e.GET("/login", ph.HandleLoginShow)
	e.GET("/signup", ph.HandleSignupShow)
	e.POST("/login", ph.HandleLogin)
	e.POST("/signup", ph.HandleSignup)
	e.GET("/logout", ph.HandleLogout, authMiddleware)

	// Admin group
	g3 := e.Group("/admin")
	g3.Use(authMiddleware, middle.Admin)
	g3.GET("", ah.HandleIndexShow)

	g3.GET("/garantia", ah.HandleGarantiaShow)
	g3.POST("/garantia", ah.HandleNewGarantiaCategory)
	g3.PUT("/garantia", ah.HandleUpdateGarantiaCategory)
	g3.DELETE("/garantia", ah.HandleRemoveGarantiaCategory)

	g3.GET("/garantia/:slug", ah.HandleGarantiaCategoryShow)
	g3.POST("/garantia/:slug", ah.HandleNewGarantiaItem)
	g3.PUT("/garantia/:slug", ah.HandleUpdateGarantiaItem)
	g3.DELETE("/garantia/:slug", ah.HandleRemoveGarantiaItem)

	g3.GET("/garantia/:categorySlug/:itemSlug", ah.HandleGarantiaItemShow)
	g3.POST("/garantia/:categorySlug/:itemSlug", ah.HandleNewGarantiaProduct)
	g3.PUT("/garantia/:categorySlug/:itemSlug", ah.HandleUpdateGarantiaProduct)
	g3.DELETE("/garantia/:categorySlug/:itemSlug", ah.HandleRemoveGarantiaProduct)

	g3.GET("/store", ah.HandleStoreShow)

	// Cart group
	g4 := e.Group("/cart")
	g4.Use(authMiddleware, cartMiddleware)
	g4.POST("", ph.HandleNewCartItem)
	g4.DELETE("", ph.HandleRemoveCartItem)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(e.Start(":" + port))
}
