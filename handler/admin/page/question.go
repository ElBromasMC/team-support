package page

import (
	"alc/handler/util"
	"alc/model/survey"
	"alc/view/admin/page/question"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleQuestionsShow(c echo.Context) error {
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
	qs, err := h.SurveyService.GetQuestions(s)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, question.Show(s, qs))
}

func (h *Handler) HandleQuestionInsertion(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}
	var q survey.Question
	q.QuestionText = c.FormValue("QuestionText")
	q.QuestionType = survey.TextType

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		return err
	}

	// Attach data
	q.Survey = s

	// Insert question
	_, err = h.SurveyService.InsertQuestion(q)
	if err != nil {
		return err
	}

	// Query updated data
	qs, err := h.SurveyService.GetQuestions(s)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, question.Table(qs))
}

func (h *Handler) HandleQuestionUpdate(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de la pregunta inválido")
	}
	var q survey.Question
	q.QuestionText = c.FormValue("QuestionText")
	q.QuestionType = survey.TextType

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		return err
	}

	// Attach data
	q.Survey = s

	// Update question
	err = h.SurveyService.UpdateQuestion(questionId, q)
	if err != nil {
		return err
	}

	// Query updated data
	qs, err := h.SurveyService.GetQuestions(s)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, question.Table(qs))
}

func (h *Handler) HandleQuestionDeletion(c echo.Context) error {
	// Parse request
	surveyId, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de encuesta inválido")
	}
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de la pregunta inválido")
	}

	// Query data
	s, err := h.SurveyService.GetSurveyById(surveyId)
	if err != nil {
		return err
	}

	// Delete question
	err = h.SurveyService.RemoveQuestion(questionId)
	if err != nil {
		return err
	}

	// Query updated data
	qs, err := h.SurveyService.GetQuestions(s)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, question.Table(qs))
}

func (h *Handler) HandleQuestionInsertionFormShow(c echo.Context) error {
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

	return util.Render(c, http.StatusOK, question.InsertionForm(s))
}

func (h *Handler) HandleQuestionUpdateFormShow(c echo.Context) error {
	// Parse request
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de la pregunta inválido")
	}

	// Query data
	q, err := h.SurveyService.GetQuestionById(questionId)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, question.UpdateForm(q))
}

func (h *Handler) HandleQuestionDeletionFormShow(c echo.Context) error {
	// Parse request
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de la pregunta inválido")
	}

	// Query data
	q, err := h.SurveyService.GetQuestionById(questionId)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, question.DeletionForm(q))
}
