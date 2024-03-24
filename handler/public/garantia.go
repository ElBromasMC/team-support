package public

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/garantia"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GET "/garantia"
func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	cats, err := h.PublicService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.Show(cats))
}

// GET "/garantia/:slug"
func (h *Handler) HandleGarantiaCategoryShow(c echo.Context) error {
	slug := c.Param("slug")

	cat, err := h.PublicService.GetCategory(store.GarantiaType, slug)
	if err != nil {
		return err
	}

	items, err := h.PublicService.GetItems(cat)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ShowCategory(cat, items))
}

// GET "/garantia/:categorySlug/:itemSlug"
func (h *Handler) HandleGarantiaItemShow(c echo.Context) error {
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	cat, err := h.PublicService.GetCategory(store.GarantiaType, categorySlug)
	if err != nil {
		return err
	}

	var item store.Item
	if err := h.PublicService.DB.QueryRow(context.Background(), `SELECT si.id, si.name, img.filename
FROM store_items AS si
LEFT JOIN images AS img
ON si.img_id = img.id
WHERE si.category_id = $1 AND si.slug = $2`, cat.Id, itemSlug).Scan(&item.Id, &item.Name, &item.Img.Filename); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	item.Slug = itemSlug
	item.Category = cat

	rows, err := h.PublicService.DB.Query(context.Background(), `SELECT id, name, price, details
FROM store_products
WHERE item_id = $1`, item.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Product, error) {
		var product store.Product
		err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Details)
		product.Item = item
		return product, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return util.Render(c, http.StatusOK, garantia.ShowItem(item, products))
}
