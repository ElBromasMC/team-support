package pages

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alc/handler/helper"

	"alc/view/pages"
)

type IndexHandler struct{}

func (h IndexHandler) HandleIndexShow(c echo.Context) error {
	return helper.Render(c, http.StatusOK, pages.Index())
}
