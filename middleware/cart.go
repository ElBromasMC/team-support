package middleware

import (
	"alc/config"
	"alc/model/cart"
	"alc/model/currency"
	"alc/service"
	"context"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Cart(ps service.Public, cs service.Currency) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get(cart.SessionName, c)
			if err != nil {
				c.Logger().Debug("Error getting cart session: ", err)
				return next(c)
			}

			// Retrieve items request
			itemsRequest, ok := sess.Values["items"].([]cart.ItemRequest)
			if !ok {
				c.Logger().Debug("Invalid items from request")
				return next(c)
			}

			// Get and normalize items from request
			items := make([]cart.Item, 0, len(itemsRequest))
			for _, i := range itemsRequest {
				item, err := ps.RequestToItem(i)
				if err != nil {
					c.Logger().Debug("Error getting item: ", err)
					return next(c)
				}
				item, err = item.Normalize()
				if err != nil {
					c.Logger().Debug("Invalid item: ", err)
					return next(c)
				}
				items = append(items, item)
			}

			// Get exchange rate
			rate, err := cs.GetExchangeRate(config.STORE_CURRENCY)
			if err != nil {
				return err
			}

			// Attach items to request context
			ctx := context.WithValue(c.Request().Context(), cart.ItemsKey{}, items)
			ctx = context.WithValue(ctx, currency.RateKey{}, rate)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
