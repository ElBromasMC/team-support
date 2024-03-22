package public

import (
	"alc/handler/util"
	"alc/view/cart"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCartShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, cart.Show())
}
