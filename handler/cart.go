package handler

import (
	"alc/view/cart"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCartShow(c echo.Context) error {
	return render(c, http.StatusOK, cart.Show())
}
