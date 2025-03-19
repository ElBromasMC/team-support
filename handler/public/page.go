package public

import (
	"alc/handler/util"
	"alc/model/book"
	"alc/model/survey"
	"alc/view/page"
	"bytes"
	"context"
	"fmt"
	"net/http"
	gomail "net/mail"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wneessen/go-mail"
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

func (h *Handler) HandleBookFormShow(c echo.Context) error {
	// Captcha verification
	captchaSiteKey := h.CaptchaService.GetSiteKey()
	return util.Render(c, http.StatusOK, page.BookForm(captchaSiteKey))
}

func (h *Handler) HandleBookEntryInsertion(c echo.Context) error {
	// Captcha verification
	// Get the reCAPTCHA token from the form.
	captchaResponse := c.FormValue("g-recaptcha-response")
	// Get the client's IP address.
	remoteIP := c.RealIP()
	// Verify the captcha token and retrieve the score.
	valid, score, err := h.CaptchaService.Verify(captchaResponse, remoteIP)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error de verificación del captcha")
	}

	// Set a threshold for the score
	if !valid || score < 0.6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Error de verificación del captcha")
	}

	// Parse request
	var entry book.Entry
	dtype, err := h.BookService.GetDocumentType(c.FormValue("DocumentType"))
	if err != nil {
		return err
	}
	gtype, err := h.BookService.GetGoodType(c.FormValue("GoodType"))
	if err != nil {
		return err
	}
	ctype, err := h.BookService.GetComplaintType(c.FormValue("ComplaintType"))
	if err != nil {
		return err
	}
	entry.DocumentType = dtype
	entry.DocumentNumber = c.FormValue("DocumentNumber")
	entry.Name = c.FormValue("Name")
	entry.Address = c.FormValue("Address")
	entry.PhoneNumber = c.FormValue("PhoneNumber")
	entry.Email = strings.ToLower(strings.TrimSpace(c.FormValue("Email")))
	entry.ParentName = c.FormValue("ParentName")
	entry.GoodType = gtype
	entry.GoodDescription = c.FormValue("GoodDescription")
	entry.ComplaintType = ctype
	entry.ComplaintDescription = c.FormValue("ComplaintDescription")
	entry.ActionsDescription = c.FormValue("ActionsDescription")

	// Validate email
	address, err := gomail.ParseAddress(entry.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Email inválido")
	}
	entry.Email = address.Address

	// Insert entry
	id, err := h.BookService.InsertBookEntry(entry)
	if err != nil {
		return err
	}

	// Query data
	entry, err = h.BookService.GetBookEntryById(id)
	if err != nil {
		return err
	}

	// Get PDF
	m, err := h.BookService.GeneratePDF(entry)
	if err != nil {
		return err
	}
	document, err := m.Generate()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	reader := bytes.NewReader(document.GetBytes())

	// Send the email
	msg := mail.NewMsg()
	if err := msg.From(h.EmailService.GetSenderEmail()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if err := msg.To(entry.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if err := msg.Bcc(h.EmailService.GetBookEmail()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	msg.Subject("Libro de Reclamaciones - Team Support Services")
	msg.SetBodyString(mail.TypeTextPlain,
		`Estimado cliente.

Se adjunta el documento generado a través del Libro de Reclamaciones.
Pronto nos pondremos en contacto con usted.

Por favor, no responder a este correo.

Saludos`,
	)
	msg.AttachReadSeeker("hoja.pdf", reader)
	if err := h.EmailService.DialAndSend(context.Background(), msg); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error en el servidor de correos")
	}

	_, ok := c.Request().Header[http.CanonicalHeaderKey("HX-Request")]
	if !ok {
		return c.Redirect(http.StatusFound, "/")
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}
