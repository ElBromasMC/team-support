package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alc/view/page"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	return render(c, http.StatusOK, page.Index())
}
