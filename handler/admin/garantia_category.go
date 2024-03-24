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

// GET "/admin/garantia/:slug"
func (h *Handler) HandleGarantiaCategoryShow(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}

	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.CategoryShow(cat, items))
}

// POST "/admin/garantia/:slug"
func (h *Handler) HandleNewGarantiaItem(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	var item store.Item
	item.Name = c.FormValue("name")
	item.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Query and attach category
	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}
	item.Category = cat

	// Generate slug from name
	item.Slug = slug.Make(item.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		item.Img = newImg
	}

	// Insert it into database
	if _, err := h.AdminService.InsertItem(item); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}

// PUT "/admin/garantia/:slug"
func (h *Handler) HandleUpdateGarantiaItem(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var item store.Item
	item.Name = c.FormValue("name")
	item.Description = c.FormValue("description")

	img, imgErr := c.FormFile("img")

	// Generate slug from name
	item.Slug = slug.Make(item.Name)

	// Query and attach category
	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}
	item.Category = cat

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		item.Img = newImg
	}

	// Update item
	if err := h.AdminService.UpdateItem(id, item); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}

// DELETE "/admin/garantia/:slug"
func (h *Handler) HandleRemoveGarantiaItem(c echo.Context) error {
	// Parsing request
	catSlug := c.Param("slug")

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	// Query category
	cat, err := h.AdminService.GetCategory(store.GarantiaType, catSlug)
	if err != nil {
		return err
	}

	// Remove item
	if err := h.AdminService.RemoveItem(id); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.ItemTableShow(items))
}
