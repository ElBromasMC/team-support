package handler

import (
	"alc/model/store"
	"alc/view/admin"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GET "/admin"
func (h *Handler) HandleAdminShow(c echo.Context) error {
	return render(c, http.StatusOK, admin.Show())
}

// GET "/admin/garantia"
func (h *Handler) HandleAdminGarantiaShow(c echo.Context) error {
	rows, err := h.DB.Query(context.Background(), `SELECT id, name, description, slug
FROM store_categories
WHERE type = $1`, store.GarantiaType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	cats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Category, error) {
		var cat store.Category
		err := row.Scan(&cat.Id, &cat.Name, &cat.Description, &cat.Slug)
		cat.Type = store.GarantiaType
		return cat, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return render(c, http.StatusOK, admin.Garantia(cats))
}

// POST "/admin/garantia"
func (h *Handler) HandleAdminGarantia(c echo.Context) error {
	return nil
}

// GET "/admin/store"
func (h *Handler) HandleAdminStoreShow(c echo.Context) error {
	return render(c, http.StatusOK, admin.Store())
}
