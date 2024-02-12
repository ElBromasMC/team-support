package middleware

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func Auth(db *pgx.Conn) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), "user", "arlr.user@gmail.com")
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
