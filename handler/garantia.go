package handler

import (
	"alc/model"
	"alc/view/garantia"
	"net/http"

	"github.com/labstack/echo/v4"
)

var items = []model.GarantiaItem{
	{Name: "Garantía + Daño Accidental + Bateria TUF", Price: 18000, Slug: "garantia-accidental-bateria", Img: "/static/img/garantia1.jpg"},
	{Name: "Protección contra daño accidental TUF", Price: 18000, Slug: "accidental", Img: "/static/img/garantia1.jpg"},
	{Name: "Garantía Extendida + Domicilio TUF", Price: 18000, Slug: "garantia-domicilio", Img: "/static/img/garantia1.jpg"},
	{Name: "Garantía Extendida TUF", Price: 18000, Slug: "garantia", Img: "/static/img/garantia1.jpg"},
}

func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	return render(c, http.StatusOK, garantia.Show(items))
}

func (h *Handler) HandleGarantiaItemShow(c echo.Context) error {
	// Item slug from path `/garantia/:slug`
	slug := c.Param("slug")

	found := false
	var item model.GarantiaItem
	for _, i := range items {
		if i.Slug == slug {
			item = i
			found = true
			break
		}
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return render(c, http.StatusOK, garantia.Item(item))
}
