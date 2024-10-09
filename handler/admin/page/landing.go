package page

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/model/survey"
	"alc/view/admin/page/landing"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleLandingsShow(c echo.Context) error {
	// Query data
	ls, err := h.SurveyService.GetLandings()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.Show(ls))
}

func (h *Handler) HandleLandingInsertion(c echo.Context) error {
	// Parse request
	var l survey.Landing
	l.Title = c.FormValue("Title")
	l.Content = c.FormValue("Content")
	surveyId, err := strconv.Atoi(c.FormValue("SurveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta no válido")
	}

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		s.Id = 0
	}

	// Attach data
	l.Survey = s

	// Insert landing
	_, err = h.SurveyService.InsertLanding(l)
	if err != nil {
		return err
	}

	// Query updated data
	ls, err := h.SurveyService.GetLandings()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.Table(ls))
}

func (h *Handler) HandleLandingUpdate(c echo.Context) error {
	// Parse request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}
	var l survey.Landing
	l.Title = c.FormValue("Title")
	l.Content = c.FormValue("Content")
	surveyId, err := strconv.Atoi(c.FormValue("SurveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta no válido")
	}

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		s.Id = 0
	}

	// Attach data
	l.Survey = s

	// Update landing
	err = h.SurveyService.UpdateLanding(landingId, l)
	if err != nil {
		return err
	}

	// Query updated data
	ls, err := h.SurveyService.GetLandings()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.Table(ls))
}

func (h *Handler) HandleLandingDeletion(c echo.Context) error {
	// Parse request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}

	// Delete landing
	err = h.SurveyService.RemoveLanding(l)
	if err != nil {
		return err
	}

	// Query updated data
	ls, err := h.SurveyService.GetLandings()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.Table(ls))
}

func (h *Handler) HandleLandingPublication(c echo.Context) error {
	// Parse request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}

	// Publish landing
	err = h.SurveyService.SetActiveLanding(l.Id)
	if err != nil {
		return err
	}

	// Query updated data
	ls, err := h.SurveyService.GetLandings()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.Table(ls))
}

func (h *Handler) HandleLandingHide(c echo.Context) error {
	// Parse request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}

	// Hide landing
	err = h.SurveyService.HideLanding(l.Id)
	if err != nil {
		return err
	}

	// Query updated data
	ls, err := h.SurveyService.GetLandings()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.Table(ls))
}

func (h *Handler) HandleLandingInsertionFormShow(c echo.Context) error {
	ss, err := h.SurveyService.GetSurveys()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.InsertionForm(ss))
}

func (h *Handler) HandleLandingUpdateFormShow(c echo.Context) error {
	// Parse request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}
	ss, err := h.SurveyService.GetSurveys()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, landing.UpdateForm(ss, l))
}

func (h *Handler) HandleLandingDeletionFormShow(c echo.Context) error {
	// Parse request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, landing.DeletionForm(l))
}

func (h *Handler) HandleLandingImagesModification(c echo.Context) error {
	// Parsing request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files, ok := form.File["imgs"]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Debe proporcionar imágenes")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}
	maxIndex, err := h.SurveyService.GetMaxIndex(l)
	if err != nil {
		return err
	}

	// Upload the images
	imgs := make([]store.Image, 0, len(files))
	for i, file := range files {
		img, err := h.SurveyService.InsertImage(file)
		if err != nil {
			continue
		}
		img.Index = (maxIndex + 1) + i
		imgs = append(imgs, img)
	}
	err = h.SurveyService.ModifyLandingImages(l, imgs)
	if err != nil {
		return err
	}

	// Get updated images
	imgs, err = h.SurveyService.GetLandingImages(l)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, landing.ModifyImagesForm(l, imgs))
}

func (h *Handler) HandleLandingImageDeletion(c echo.Context) error {
	// Parsing request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}
	imgId, err := strconv.Atoi(c.FormValue("Id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id no válido")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}

	// Delete image
	if err := h.SurveyService.RemoveImage(imgId); err != nil {
		return err
	}

	// Get updated images
	imgs, err := h.SurveyService.GetLandingImages(l)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, landing.ModifyImagesForm(l, imgs))
}

func (h *Handler) HandleLandingImagesFormShow(c echo.Context) error {
	// Parse request
	landingId, err := strconv.Atoi(c.Param("landingId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de landing inválido")
	}

	// Query data
	l, err := h.SurveyService.GetLandingById(landingId)
	if err != nil {
		return err
	}
	imgs, err := h.SurveyService.GetLandingImages(l)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, landing.ModifyImagesForm(l, imgs))
}
