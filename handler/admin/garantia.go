package admin

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/admin/garantia"
	"context"
	"net/http"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

// GET "/admin/garantia"
func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	cats, err := h.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.Show(cats))
}

// POST "/admin/garantia"
func (h *Handler) HandleNewGarantiaCategory(c echo.Context) error {
	name := c.FormValue("name")
	description := c.FormValue("description")
	slug := slug.Make(name)

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
	return util.Render(c, http.StatusOK, garantia.CategoryTableShow(cats))
}

// PUT "/admin/garantia"
func (h *Handler) HandleUpdateGarantiaCategory(c echo.Context) error {
	idStr := c.FormValue("id")

	name := c.FormValue("name")
	description := c.FormValue("description")
	slug := slug.Make(name)
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
		var prevImgId *int
		var prevImgFilename *string
		h.DB.QueryRow(context.Background(), `SELECT img.id, img.filename
FROM store_categories AS sc
LEFT JOIN images AS img
ON sc.img_id = img.id
WHERE sc.id = $1`, id).Scan(&prevImgId, &prevImgFilename)
		// Remove previous image if exists
		if prevImgId != nil && prevImgFilename != nil {
			h.RemoveImage(store.Image{
				Id:       *prevImgId,
				Filename: *prevImgFilename,
			})
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
	return util.Render(c, http.StatusOK, garantia.CategoryTableShow(cats))
}

// GET "/admin/garantia/:slug"
func (h *Handler) HandleGarantiaCategoryShow(c echo.Context) error {
	slug := c.Param("slug")
	cat, err := h.GetCategory(store.GarantiaType, slug)
	if err != nil {
		return err
	}

	items, err := h.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryShow(cat, items))
}

// POST "/admin/garantia/:slug"
func (h *Handler) HandleNewGarantiaItem(c echo.Context) error {
	catSlug := c.Param("slug")
	name := c.FormValue("name")
	description := c.FormValue("description")
	slug := slug.Make(name)

	cat, err := h.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}

	img, err := c.FormFile("img")
	if err != nil {
		if _, err := h.DB.Exec(context.Background(), `INSERT INTO store_items (category_id, name, description, slug)
VALUES ($1, $2, $3, $4)`, cat.Id, name, description, slug); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new item into database")
		}
	} else {
		imgId, err := h.InsertImage(img)
		if err != nil {
			return err
		}

		if _, err := h.DB.Exec(context.Background(), `INSERT INTO store_items (category_id, name, description, img_id, slug)
VALUES ($1, $2, $3, $4, $5)`, cat.Id, name, description, imgId, slug); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new item into database")
		}
	}

	items, err := h.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}

// PUT "/admin/garantia/:slug"
func (h *Handler) HandleUpdateGarantiaItem(c echo.Context) error {
	catSlug := c.Param("slug")
	idStr := c.FormValue("id")

	name := c.FormValue("name")
	description := c.FormValue("description")
	slug := slug.Make(name)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	img, err := c.FormFile("img")
	if err != nil {
		if _, err := h.DB.Exec(context.Background(), `UPDATE store_items
SET name = $1, description = $2, slug = $3
WHERE id = $4`, name, description, slug, id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error updating the item into database")
		}
	} else {
		var prevImgId *int
		var prevImgFilename *string
		h.DB.QueryRow(context.Background(), `SELECT img.id, img.filename
FROM store_items AS si
LEFT JOIN images AS img
ON si.img_id = img.id
WHERE si.id = $1`, id).Scan(&prevImgId, &prevImgFilename)
		// Remove previous image if exists
		if prevImgId != nil && prevImgFilename != nil {
			h.RemoveImage(store.Image{
				Id:       *prevImgId,
				Filename: *prevImgFilename,
			})
		}

		// Insert new image
		imgId, err := h.InsertImage(img)
		if err != nil {
			return err
		}

		if _, err := h.DB.Exec(context.Background(), `UPDATE store_items
SET name = $1, description = $2, img_id = $3, slug = $4
WHERE id = $5`, name, description, imgId, slug, id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error updating the item into database")
		}
	}

	cat, err := h.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}

	items, err := h.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}
