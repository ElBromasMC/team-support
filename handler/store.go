package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alc/model"

	"alc/view/store"
)

var item0 model.StoreItem
var storeItems []model.StoreItem

func init() {
	item0 = model.StoreItem{
		Category:         "Cargadores",
		Name:             "SERVIDOR XEON 2224G DELL PowerEdge T40 8|1TB Torre",
		Price:            100,
		BriefDescription: "El componente fundamental para su pequeña empresa. PowerEdge T40, confiable, eficiente y asequible para brindar soporte a su empresa en crecimiento. Ideal para archivos, impresiones, correos y mensajería.",
		Description: `Procesador:
Intel Xeon E-2224G 3.5GHZ (caché de 8 M, 4.70 GHz)
Memoria Ram:
8GB DDR4 2666 (1x8GB) 4 Ranuras
Disco Duro:
1TB SATA 7200rpm`,
		Slug:     "servidor-xeon-poweredge-1tb",
		Img:      "/static/img/store1.jpg",
		LargeImg: "/static/img/_store1.jpg",
	}
	for i := 0; i < 9; i++ {
		storeItems = append(storeItems, item0)
	}
}

func (h *Handler) HandleStoreShow(c echo.Context) error {
	return render(c, http.StatusOK, store.Show(storeItems))
}

func (h *Handler) HandleStoreItemShow(c echo.Context) error {
	// Item slug from path `/store/:slug`
	slug := c.Param("slug")

	var item model.StoreItem
	found := false
	for _, i := range storeItems {
		if i.Slug == slug {
			item = i
			found = true
			break
		}
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return render(c, http.StatusOK, store.ShowItem(item))
}
