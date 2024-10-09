package public

import (
	"alc/handler/util"
	"alc/model/survey"
	"alc/view/page"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Index())
}

func (h *Handler) HandleTicketShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Ticket())
}

func (h *Handler) HandleLandingShow(c echo.Context) error {
	// Parse request
	l, err := h.SurveyService.GetActiveLanding()
	if err != nil {
		return err
	}
	var qs []survey.Question
	if l.Survey.Id != 0 {
		qs, err = h.SurveyService.GetQuestions(l.Survey)
		if err != nil {
			return err
		}
	}
	imgs, err := h.SurveyService.GetLandingImages(l)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, page.Landing(l, qs, imgs))
}

func (h *Handler) HandleSurveyInsertion(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuestas inválido")
	}
	var response survey.SurveyResponse
	response.Name = c.FormValue("Name")
	response.PhoneNumber = c.FormValue("PhoneNumber")
	response.Email = c.FormValue("Email")
	rating, err := strconv.Atoi(c.FormValue("Rating"))
	// Validation
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Rating inválido")
	}
	if !(1 <= rating && rating <= 5) {
		return echo.NewHTTPError(http.StatusBadRequest, "El rating debe ser un valor del 1 al 5")
	}
	response.Rating = rating

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		return err
	}
	qs, err := h.SurveyService.GetQuestions(s)
	if err != nil {
		return err
	}

	// Attach data
	response.QuestionResponses = make([]survey.QuestionResponse, 0, len(qs))
	for _, q := range qs {
		response.QuestionResponses = append(response.QuestionResponses, survey.QuestionResponse{
			Question:     q,
			ResponseText: c.FormValue(fmt.Sprintf("Question_%d", q.Id)),
		})
	}
	response.Survey = s
	err = h.SurveyService.InsertResponse(response)
	if err != nil {
		return err
	}

	return nil
}
