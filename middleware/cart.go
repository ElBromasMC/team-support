package middleware

import (
	"alc/handler/util"
	"alc/model/cart"
	"alc/service"
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func Cart(ps service.Public) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Read cart cookie
			cookie, err := c.Cookie(cart.CookieName)
			if err != nil {
				util.RemoveCookie(c, cart.CookieName)
				return next(c)
			}

			// Parse item requests
			var itemRequests []cart.ItemRequest
			if err := json.Unmarshal([]byte(cookie.Value), &itemRequests); err != nil {
				util.RemoveCookie(c, cart.CookieName)
				return next(c)
			}

			// Validate data
			items := make([]cart.Item, 0, len(itemRequests))
			for _, i := range itemRequests {
				item, err := i.ToValidItem(ps)
				if err != nil {
					continue
				}
				items = append(items, item)
			}

			// Overwrite cart cookie
			if err := cart.PutCookie(c, items); err != nil {
				util.RemoveCookie(c, cart.CookieName)
				return next(c)
			}

			ctx := context.WithValue(c.Request().Context(), cart.ItemsKey{}, items)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
