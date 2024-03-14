package middleware

import (
	"alc/handler"
	"alc/model"
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Auth(db *pgxpool.Pool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return next(c)
			}
			session, err := uuid.FromString(cookie.Value)
			if err != nil {
				handler.RemoveCookie(c, "session")
				return next(c)
			}
			var user model.User
			if err = db.QueryRow(context.Background(), `SELECT u.name, u.email, u.role
FROM users u
JOIN sessions s ON u.user_id = s.user_id
WHERE s.session_id = $1`, session).Scan(&user.Name, &user.Email, &user.Role); err != nil {
				handler.RemoveCookie(c, "session")
				return next(c)
			}
			user.Session = session
			ctx := context.WithValue(c.Request().Context(), "user", user)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
