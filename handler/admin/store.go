package admin

import (
	"alc/handler/util"
	"alc/view/admin/store"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET "/admin/store"
func (h *Handler) HandleStoreShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, store.Show())
}
