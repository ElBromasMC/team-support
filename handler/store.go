package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alc/model"

	"alc/view/store"
)

func (h *Handler) HandleStoreShow(c echo.Context) error {
	items := []model.StoreItem{
		{
			Name:  "hola",
			Price: 200,
		},
	}
	return render(c, http.StatusOK, store.Show(items))
}
