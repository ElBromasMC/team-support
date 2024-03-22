package public

import (
	"alc/model/store"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB *pgxpool.Pool
}

func (h *Handler) GetCategories(t store.Type) ([]store.Category, error) {
	rows, err := h.DB.Query(context.Background(), `SELECT sc.id, sc.name, sc.description, sc.slug, img.filename
FROM store_categories AS sc
LEFT JOIN images AS img
ON sc.img_id = img.id
WHERE type = $1
ORDER BY id DESC`, t)
	if err != nil {
		return []store.Category{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	cats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Category, error) {
		var cat store.Category
		var img *string
		err := row.Scan(&cat.Id, &cat.Name, &cat.Description, &cat.Slug, &img)
		if img != nil {
			cat.Img.Filename = *img
		} else {
			cat.Img.Filename = ""
		}
		cat.Type = t
		return cat, err
	})
	if err != nil {
		return []store.Category{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return cats, nil
}

func (h *Handler) GetItems(cat store.Category) ([]store.Item, error) {
	rows, err := h.DB.Query(context.Background(), `SELECT si.id, si.name, si.description, si.long_description, si.slug, img.filename, largeimg.filename
FROM store_items AS si
LEFT JOIN images AS img
ON si.img_id = img.id
LEFT JOIN images as largeimg
ON si.largeimg_id = largeimg.id
WHERE si.category_id = $1
ORDER BY id DESC`, cat.Id)
	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Item, error) {
		var item store.Item
		var img *string
		var largeimg *string
		err := row.Scan(&item.Id, &item.Name, &item.Description, &item.LongDescription, &item.Slug, &img, &largeimg)
		if img != nil {
			item.Img.Filename = *img
		} else {
			item.Img.Filename = ""
		}
		if largeimg != nil {
			item.LargeImg.Filename = *largeimg
		} else {
			item.LargeImg.Filename = ""
		}
		item.Category = cat
		return item, err
	})
	if err != nil {
		return []store.Item{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return items, nil
}
