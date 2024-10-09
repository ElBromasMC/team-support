package page

import (
	"alc/handler/util"
	"alc/model/survey"
	view "alc/view/admin/page/survey"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleSurveysShow(c echo.Context) error {
	// Query data
	ss, err := h.SurveyService.GetSurveys()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, view.Show(ss))
}

func (h *Handler) HandleSurveyInsertion(c echo.Context) error {
	// Parse request
	var s survey.Survey
	s.Title = c.FormValue("Title")

	// Insert survey
	_, err := h.SurveyService.InsertSurvey(s)
	if err != nil {
		return err
	}

	// Query updated data
	ss, err := h.SurveyService.GetSurveys()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, view.Table(ss))
}

func (h *Handler) HandleSurveyUpdate(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}
	var s survey.Survey
	s.Title = c.FormValue("Title")

	// Update survey
	err = h.SurveyService.UpdateSurvey(surveyId, s)
	if err != nil {
		return err
	}

	// Query updated data
	ss, err := h.SurveyService.GetSurveys()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, view.Table(ss))
}

func (h *Handler) HandleSurveyDeletion(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}

	// Delete survey
	err = h.SurveyService.RemoveSurvey(surveyId)
	if err != nil {
		return err
	}

	// Query updated data
	ss, err := h.SurveyService.GetSurveys()
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, view.Table(ss))
}

func (h *Handler) HandleSurveyResultsDownload(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		return err
	}

	// Write the csv
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=\"data.csv\"")

	err = h.SurveyService.WriteResponses(c.Response().Writer, s)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleSurveyInsertionFormShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, view.InsertionForm())
}

func (h *Handler) HandleSurveyUpdateFormShow(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, view.UpdateForm(s))
}

func (h *Handler) HandleSurveyDeletionFormShow(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, view.DeletionForm(s))
}
