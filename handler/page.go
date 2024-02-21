package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alc/model"
	"alc/view/page"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	items := []model.ServiceItem{
		{Name: "Soporte para hardware", Description: "Cambio y reparación de partes para tu equipo"},
		{Name: "Renting de equipos", Description: "Alquiler de equipo con soporte por contrato"},
		{Name: "Proyectos de T.I.", Description: "Te ayudamos a planificar y realizar tus proyectos de T.I."},
		{Name: "Garantías", Description: "Seguro de vida para tus equipos"},
	}

	return render(c, http.StatusOK, page.Index(items))
}

func (h *Handler) HandleTicketShow(c echo.Context) error {
	return render(c, http.StatusOK, page.Ticket())
}

func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	items := []model.GarantiaItem{
		{Name: "Garantía + Daño Accidental + Bateria TUF", Price: 18000, Slug: "garantia-accidental-bateria", Img: "/static/img/garantia1.jpg"},
		{Name: "Protección contra daño accidental TUF", Price: 18000, Slug: "accidental", Img: "/static/img/garantia1.jpg"},
		{Name: "Garantía Extendida + Domicilio TUF", Price: 18000, Slug: "garantia-domicilio", Img: "/static/img/garantia1.jpg"},
		{Name: "Garantía Extendida TUF", Price: 18000, Slug: "garantia", Img: "/static/img/garantia1.jpg"},
	}
	return render(c, http.StatusOK, page.Garantia(items))
}

func (h *Handler) HandleGarantiaItemShow(c echo.Context) error {
	// Item slug from path `/garantia/:slug`
	slug := c.Param("slug")

	items := []model.GarantiaItem{
		{Name: "Garantía + Daño Accidental + Bateria TUF", Price: 18923, Slug: "garantia-accidental-bateria", Img: "/static/img/garantia1.jpg"},
		{Name: "Protección contra daño accidental TUF", Price: 18000, Slug: "accidental", Img: "/static/img/garantia1.jpg"},
		{Name: "Garantía Extendida + Domicilio TUF", Price: 18000, Slug: "garantia-domicilio", Img: "/static/img/garantia1.jpg"},
		{Name: "Garantía Extendida TUF", Price: 18000, Slug: "garantia", Img: "/static/img/garantia1.jpg"},
	}

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

	return render(c, http.StatusOK, page.GarantiaItem(item))
}
