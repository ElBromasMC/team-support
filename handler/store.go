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
			Category: "Cargadores",
			Name:     "SERVIDOR XEON 2224G DELL PowerEdge T40 8|1TB Torre",
			Price:    300,
			Slug:     "servidor-xeon-poweredge-1tb",
			Img:      "/static/img/store1.jpg",
		},
	}
	return render(c, http.StatusOK, store.Show(items))
}
