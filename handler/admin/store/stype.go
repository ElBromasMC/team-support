package store

import (
	"alc/handler/util"
	"alc/view/admin/store/stype"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleStoreTypeShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Tipo inválido")
	}

	return util.Render(c, http.StatusOK, stype.Show(t))
}
