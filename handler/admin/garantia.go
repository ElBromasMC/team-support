package admin

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/admin/garantia"
	"net/http"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

// GET "/admin/garantia"
func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	cats, err := h.AdminService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.Show(cats))
}

// POST "/admin/garantia"
func (h *Handler) HandleNewGarantiaCategory(c echo.Context) error {
	// Parsing request
	var cat store.Category
	cat.Type = store.GarantiaType
	cat.Name = c.FormValue("name")
	cat.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Generate slug from name
	cat.Slug = slug.Make(cat.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		cat.Img = newImg
	}

	// Insert it into database
	if _, err := h.AdminService.InsertCategory(cat); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryTableShow(cats))
}

// PUT "/admin/garantia"
func (h *Handler) HandleUpdateGarantiaCategory(c echo.Context) error {
	// Parsing request
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var cat store.Category
	cat.Type = store.GarantiaType
	cat.Name = c.FormValue("name")
	cat.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Generate slug from name
	cat.Slug = slug.Make(cat.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		cat.Img = newImg
	}

	// Update category
	if err := h.AdminService.UpdateCategory(id, cat); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryTableShow(cats))
}

// DELETE "/admin/garantia"
func (h *Handler) HandleRemoveGarantiaCategory(c echo.Context) error {
	// Parsing request
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	// Remove category
	if err := h.AdminService.RemoveCategory(id); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryTableShow(cats))
}
