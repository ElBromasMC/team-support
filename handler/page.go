package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alc/model"
	"alc/view/page"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	items := []model.ServiceItem{
		{
			Name:            "Soporte para hardware",
			Description:     "Cambio y reparación de partes para tu equipo",
			LongDescription: "Cambio y reparación de partes para tu equipo, contamos con certificación de empresa reconocidas para el cambio de partes originales.",
		},
		{
			Name:            "Renting de equipos",
			Description:     "Alquiler de equipo con soporte por contrato",
			LongDescription: "Cambio y reparación de partes para tu equipo, contamos con certificación de empresa reconocidas para el cambio de partes originales.",
		},
		{
			Name:            "Proyectos de T.I.",
			Description:     "Te ayudamos a planificar y realizar tus proyectos de T.I.",
			LongDescription: "Cambio y reparación de partes para tu equipo, contamos con certificación de empresa reconocidas para el cambio de partes originales.",
		},
		{
			Name:            "Garantías",
			Description:     "Seguro de vida para tus equipos",
			LongDescription: "Cambio y reparación de partes para tu equipo, contamos con certificación de empresa reconocidas para el cambio de partes originales.",
		},
	}

	return render(c, http.StatusOK, page.Index(items))
}

func (h *Handler) HandleTicketShow(c echo.Context) error {
	return render(c, http.StatusOK, page.Ticket())
}
