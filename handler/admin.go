package handler

import (
	"alc/config"
	"alc/model/store"
	"alc/view/admin"
	"context"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GET "/admin"
func (h *Handler) HandleAdminShow(c echo.Context) error {
	return render(c, http.StatusOK, admin.Show())
}

// GET "/admin/garantia"
func (h *Handler) HandleAdminGarantiaShow(c echo.Context) error {
	cats, err := h.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return render(c, http.StatusOK, admin.CategoryShow(cats))
}

// POST "/admin/garantia"
func (h *Handler) HandleNewGarantia(c echo.Context) error {
	name := c.FormValue("name")
	description := c.FormValue("description")
	slug := c.FormValue("slug")

	img, err := c.FormFile("img")
	if err != nil {
		if _, err := h.DB.Exec(context.Background(), `INSERT INTO store_categories (type, name, description, slug)
VALUES ($1, $2, $3, $4)`, store.GarantiaType, name, description, slug); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new category into database")
		}
	} else {
		imgId, err := h.InsertImage(img)
		if err != nil {
			return err
		}

		if _, err := h.DB.Exec(context.Background(), `INSERT INTO store_categories (type, name, description, img_id, slug)
VALUES ($1, $2, $3, $4, $5)`, store.GarantiaType, name, description, imgId, slug); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new category into database")
		}
	}

	cats, err := h.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return render(c, http.StatusOK, admin.CategoryTableShow(cats))
}

// PUT "/admin/garantia"
func (h *Handler) HandleUpdateGarantia(c echo.Context) error {
	idStr := c.FormValue("id")
	name := c.FormValue("name")
	description := c.FormValue("description")
	slug := c.FormValue("slug")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	img, err := c.FormFile("img")
	if err != nil {
		if _, err := h.DB.Exec(context.Background(), `UPDATE store_categories
SET name = $1, description = $2, slug = $3
WHERE id = $4`, name, description, slug, id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error updating the category into database")
		}
	} else {
		// Remove previous image if exists
		var prevImgId *int
		var prevImgFilename *string
		if err := h.DB.QueryRow(context.Background(), `SELECT img.id, img.filename
FROM store_categories AS sc
LEFT JOIN images AS img
ON sc.img_id = img.id
WHERE sc.id = $1`, id).Scan(&prevImgId, &prevImgFilename); err == nil {
			// Delete from filesystem
			os.Remove(path.Join(config.IMAGES_SAVEDIR, *prevImgFilename))
			// Delete from database
			h.DB.Exec(context.Background(), `DELETE FROM images WHERE id = $1`, *prevImgId)
		}

		// Insert new image
		imgId, err := h.InsertImage(img)
		if err != nil {
			return err
		}

		if _, err := h.DB.Exec(context.Background(), `UPDATE store_categories
SET name = $1, description = $2, img_id = $3, slug = $4
WHERE id = $5`, name, description, imgId, slug, id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error updating the category into database")
		}
	}

	cats, err := h.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return render(c, http.StatusOK, admin.CategoryTableShow(cats))
}

// GET "/admin/store"
func (h *Handler) HandleAdminStoreShow(c echo.Context) error {
	return render(c, http.StatusOK, admin.Store())
}
