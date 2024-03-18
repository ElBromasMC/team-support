package handler

import (
	"alc/model/store"
	view "alc/view/store"
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GET "/store/all"
func (h *Handler) HandleStoreAllShow(c echo.Context) error {
	rows, err := h.DB.Query(context.Background(), `SELECT name, slug
FROM store_categories
WHERE type = $1`, store.StoreType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	cats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Category, error) {
		var cat store.Category
		err := row.Scan(&cat.Name, &cat.Slug)
		cat.Type = store.StoreType
		return cat, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return render(c, http.StatusOK, view.Show(cats, "all"))
}

// GET "/store/:slug"
func (h *Handler) HandleStoreCategoryShow(c echo.Context) error {
	slug := c.Param("slug")

	rows, err := h.DB.Query(context.Background(), `SELECT name, slug
FROM store_categories
WHERE sc.type = $1`, store.StoreType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	found := false
	cats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Category, error) {
		var cat store.Category
		err := row.Scan(&cat.Name, &cat.Slug)
		cat.Type = store.StoreType
		if cat.Slug == slug {
			found = true
		}
		return cat, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return render(c, http.StatusOK, view.Show(cats, slug))
}

// POST "/store/all"
func (h *Handler) HandleStoreAllItemsShow(c echo.Context) error {
	like := c.QueryParam("like")
	p := c.QueryParam("p")

	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		page = 1
	}

	rows, err := h.DB.Query(context.Background(), `SELECT si.name, si.slug, sc.slug, img.filename
FROM store_items AS si
JOIN store_categories AS sc
ON si.category_id = sc.id
LEFT JOIN images AS img
ON si.img_id = img.id
WHERE sc.type = $1
AND si.name % $2
ORDER BY si.name <-> $2
LIMIT 10 OFFSET ($3 - 1) * 9`, store.StoreType, like, page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Item, error) {
		var item store.Item
		err := row.Scan(&item.Name, &item.Slug, &item.Category.Slug, &item.Img.Filename)
		return item, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return render(c, http.StatusOK, view.ShowItems(items))
}

// POST "/store/:slug"
func (h *Handler) HandleStoreCategoryItemsShow(c echo.Context) error {
	slug := c.Param("slug")
	like := c.QueryParam("like")
	p := c.QueryParam("p")

	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		page = 1
	}

	rows, err := h.DB.Query(context.Background(), `SELECT si.name, si.slug, sc.slug, img.filename
FROM store_items AS si
JOIN store_categories AS sc
ON si.category_id = sc.id
LEFT JOIN images AS img
ON si.img_id = img.id
WHERE sc.type = $1 AND sc.slug = $2
AND si.name % $3
ORDER BY si.name <-> $3
LIMIT 10 OFFSET ($4 - 1) * 9`, store.StoreType, slug, like, page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Item, error) {
		var item store.Item
		err := row.Scan(&item.Name, &item.Slug, &item.Category.Slug, &item.Img.Filename)
		return item, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return render(c, http.StatusOK, view.ShowItems(items))
}

// GET "/store/:categorySlug/:itemSlug"
func (h *Handler) HandleStoreItemShow(c echo.Context) error {
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	var cat store.Category
	if err := h.DB.QueryRow(context.Background(), `SELECT id, name
FROM store_categories
WHERE type = $1 AND slug = $2`, store.StoreType, categorySlug).Scan(&cat.Id, &cat.Name); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	cat.Type = store.StoreType
	cat.Slug = categorySlug

	var item store.Item
	if err := h.DB.QueryRow(context.Background(), `SELECT si.id, si.name, si.description, si.long_description, largeimg.filename
FROM store_items AS si
LEFT JOIN images AS largeimg
ON si.largeimg_id = largeimg.id
WHERE si.category_id = $1
AND si.slug = $2`, cat.Id, itemSlug).Scan(&item.Id, &item.Name, &item.Description, &item.LongDescription, &item.LargeImg.Filename); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	item.Slug = itemSlug
	item.Category = cat

	rows, err := h.DB.Query(context.Background(), `SELECT id, name, price, details
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

	return render(c, http.StatusOK, view.ShowItem(item, products))
}
