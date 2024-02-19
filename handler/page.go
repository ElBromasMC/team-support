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
