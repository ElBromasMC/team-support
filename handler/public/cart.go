package public

import (
	"alc/handler/util"
	"alc/model/cart"
	"alc/model/store"
	view "alc/view/cart"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// POST "/cart"
func (h *Handler) HandleNewCartItem(c echo.Context) error {
	// Bind
	var itemRequest cart.ItemRequest
	id, err := strconv.Atoi(c.FormValue("Id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}
	itemRequest.Id = id
	quantity, err := strconv.Atoi(c.FormValue("Quantity"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Quantity")
	}
	itemRequest.Quantity = quantity
	serie := c.FormValue("Details[Serie]")
	if len(serie) != 0 {
		itemRequest.Details = make(map[string]string)
		itemRequest.Details["Serie"] = serie
	}

	// Validate data
	newItem, err := itemRequest.ToValidItem(h.PublicService)
	if err != nil {
		return err
	}

	// Get cart items
	items := cart.GetItems(c.Request().Context())

	// Insert new cart item
	found := false
	for n, i := range items {
		if i.Product.Item.Id != newItem.Product.Item.Id {
			continue
		}
		if i.Product.Item.Category.Type != store.StoreType {
			if !strings.EqualFold(i.Details["Serie"], newItem.Details["Serie"]) {
				continue
			}
			found = true
			items[n] = newItem
			continue
		}
		if i.Product.Id != newItem.Product.Id {
			continue
		}
		found = true
		items[n].Quantity += newItem.Quantity
	}
	if !found {
		items = append(items, newItem)
	}

	// Overwrite cart cookie
	if err := cart.PutCookie(c, items); err != nil {
		util.RemoveCookie(c, cart.CookieName)
		return err
	}

	return util.Render(c, http.StatusOK, view.Show(items))
}
