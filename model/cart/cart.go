package cart

import (
	"alc/model/store"
	"alc/service"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const CookieName = "cart"

type ItemsKey struct{}

type Item struct {
	Product  store.Product
	Quantity int
	Details  map[string]string
}

type ItemRequest struct {
	Id       int
	Quantity int
	Details  map[string]string
}

type OrderRequest struct {
	Order        store.Order
	ItemRequests []ItemRequest
}

func (item Item) IsValid() error {
	if item.Product.Item.Category.Type == store.GarantiaType {
		if item.Quantity != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity for warranty")
		}
		serie, ok := item.Details["Serie"]
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing 'Serie' for warranty")
		}
		if !(12 <= len(serie) && len(serie) <= 15) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid 'Serie' for warranty")
		}
	} else {
		if item.Quantity < 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity for store item")
		}
	}
	return nil
}

func (i Item) ToRequest() ItemRequest {
	return ItemRequest{
		Id:       i.Product.Id,
		Quantity: i.Quantity,
		Details:  i.Details,
	}
}

func (i ItemRequest) ToValidItem(ps service.Public) (Item, error) {
	product, err := ps.GetProductById(i.Id)
	if err != nil {
		return Item{}, err
	}
	item := Item{
		Product:  product,
		Quantity: i.Quantity,
		Details:  i.Details,
	}
	if err := item.IsValid(); err != nil {
		return Item{}, err
	}
	return item, nil
}

func PutCookie(c echo.Context, items []Item) error {
	itemRequests := make([]ItemRequest, 0, len(items))
	for _, i := range items {
		itemRequests = append(itemRequests, i.ToRequest())
	}

	encoded, err := json.Marshal(itemRequests)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	cookie := new(http.Cookie)
	cookie.Name = CookieName
	cookie.Value = string(encoded)
	cookie.Expires = time.Now().AddDate(0, 0, 1)
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
	return nil
}

func GetItems(ctx context.Context) []Item {
	items, _ := ctx.Value(ItemsKey{}).([]Item)
	return items
}
