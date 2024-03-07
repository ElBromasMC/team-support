package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alc/view/page"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	return render(c, http.StatusOK, page.Index())
}

func (h *Handler) HandleTicketShow(c echo.Context) error {
	return render(c, http.StatusOK, page.Ticket())
}
