package handler

import (
	"alc/view/admin"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleAdminShow(c echo.Context) error {
	return render(c, http.StatusOK, admin.Show())
}
