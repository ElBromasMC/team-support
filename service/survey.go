package service

import (
	"alc/config"
	"alc/model/store"
	"alc/model/survey"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Survey struct {
	db *pgxpool.Pool
}

func NewSurveyService(db *pgxpool.Pool) Survey {
	return Survey{
		db: db,
	}
}

// Survey management

func (ss Survey) GetSurveys() ([]survey.Survey, error) {
	sql := `
	SELECT
		id,
		title,
		created_at
	FROM
		surveys
	ORDER BY
		created_at DESC
	`
	rows, err := ss.db.Query(context.Background(), sql)
	if err != nil {
		return []survey.Survey{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()
	surveys, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (survey.Survey, error) {
		var s survey.Survey
		err := row.Scan(
			&s.Id,
			&s.Title,
			&s.CreatedAt,
		)
		return s, err
	})
	if err != nil {
		return []survey.Survey{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return surveys, nil
}

func (ss Survey) GetSurveyById(id int) (survey.Survey, error) {
	sql := `
	SELECT
		id,
		title,
		created_at
	FROM
		surveys
	WHERE
		id = $1
	`
	var s survey.Survey
	err := ss.db.QueryRow(context.Background(), sql, id).Scan(
		&s.Id,
		&s.Title,
		&s.CreatedAt,
	)
	if err != nil {
		return survey.Survey{}, echo.NewHTTPError(http.StatusNotFound, "Encuesta no encontrada")
	}
	return s, nil
}

func (ss Survey) InsertSurvey(s survey.Survey) (int, error) {
	sql := `
	INSERT INTO surveys (
		title
	)
	VALUES (
		$1
	)
	RETURNING
		id
	`
	var surveyId int
	err := ss.db.QueryRow(context.Background(), sql, s.Title).Scan(&surveyId)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return surveyId, nil
}

func (ss Survey) UpdateSurvey(id int, uptSurvey survey.Survey) error {
	sql := `
	UPDATE surveys
	SET
		title = $2
	WHERE
		id = $1
	`
	_, err := ss.db.Exec(context.Background(), sql, id, uptSurvey.Title)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func (ss Survey) RemoveSurvey(id int) error {
	sql := `DELETE FROM surveys WHERE id = $1`
	_, err := ss.db.Exec(context.Background(), sql, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// Question management

func (ss Survey) GetQuestionType(slug string) (survey.QuestionType, error) {
	if slug == "TEXT" {
		return survey.TextType, nil
	} else {
		return "", echo.NewHTTPError(http.StatusNotFound, "Tipo de pregunta inválido")
	}
}

func (ss Survey) GetQuestions(s survey.Survey) ([]survey.Question, error) {
	sql := `
	SELECT
		id,
		question_text,
		question_type
	FROM
		survey_questions
	WHERE
		survey_id = $1
	ORDER BY
		id
	`
	rows, err := ss.db.Query(context.Background(), sql, s.Id)
	if err != nil {
		return []survey.Question{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()
	questions, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (survey.Question, error) {
		var question survey.Question
		err := row.Scan(&question.Id, &question.QuestionText, &question.QuestionType)
		question.Survey = s
		return question, err
	})
	if err != nil {
		return []survey.Question{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return questions, nil
}

// TODO*
func (ss Survey) GetQuestionById(id int) (question survey.Question, err error) {
	sql := `
	SELECT
		id,
		question_text,
		question_type,
		survey_id
	FROM
		survey_questions
	WHERE
		id = $1
	`
	err = ss.db.QueryRow(context.Background(), sql, id).Scan(
		&question.Id,
		&question.QuestionText,
		&question.QuestionType,
		&question.Survey.Id,
	)
	if err != nil {
		err = echo.NewHTTPError(http.StatusNotFound, "Pregunta no encontrada")
		return
	}
	s, err := ss.GetSurveyById(question.Survey.Id)
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	question.Survey = s
	return
}

func (ss Survey) InsertQuestion(question survey.Question) (int, error) {
	sql := `
	INSERT INTO survey_questions (
		survey_id,
		question_text,
		question_type
	)
	VALUES (
		$1,
		$2,
		$3
	)
	RETURNING
		id
	`
	var questionId int
	err := ss.db.QueryRow(context.Background(), sql,
		question.Survey.Id,
		question.QuestionText,
		question.QuestionType,
	).Scan(&questionId)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return questionId, nil
}

func (ss Survey) UpdateQuestion(id int, uptQuestion survey.Question) error {
	sql := `
	UPDATE survey_questions
	SET
		question_text = $2
	WHERE
		id = $1
	`
	_, err := ss.db.Exec(context.Background(), sql, id, uptQuestion.QuestionText)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func (ss Survey) RemoveQuestion(id int) error {
	sql := `DELETE FROM survey_questions WHERE id = $1`
	_, err := ss.db.Exec(context.Background(), sql, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// Responses management

func (ss Survey) GetResponses(s survey.Survey) ([]survey.Question, []survey.SurveyResponse, error) {
	qs, err := ss.GetQuestions(s)
	if err != nil {
		return []survey.Question{}, []survey.SurveyResponse{}, err
	}

	if len(qs) > 0 {
		q_ids := make([]int, 0, len(qs))
		for _, q := range qs {
			q_ids = append(q_ids, q.Id)
		}

		sql := `
		WITH question_list AS (
			SELECT 
				unnest($1::int[]) AS question_id,
				generate_series(1, array_length($1::int[], 1)) AS idx
		)
		SELECT
			sr.id,
			sr.name,
			sr.email,
			sr.phone_number,
			sr.rating,
			sr.created_at,
			array_agg(
				COALESCE(qr.response_text, '') ORDER BY ql.idx
			) AS responses
		FROM
			survey_respondents sr
			CROSS JOIN question_list ql
			LEFT JOIN question_responses qr ON sr.id = qr.respondent_id AND ql.question_id = qr.question_id
		WHERE
			sr.survey_id = $2
		GROUP BY
			sr.id, sr.name, sr.email, sr.phone_number, sr.rating, sr.created_at
		ORDER BY
			sr.created_at
		`
		rows, err := ss.db.Query(context.Background(), sql, q_ids, s.Id)
		if err != nil {
			return []survey.Question{}, []survey.SurveyResponse{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
		defer rows.Close()

		responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (survey.SurveyResponse, error) {
			var response survey.SurveyResponse
			err := row.Scan(
				&response.Id,
				&response.Name,
				&response.Email,
				&response.PhoneNumber,
				&response.Rating,
				&response.CreatedAt,
				&response.Responses,
			)
			response.Survey = s
			return response, err
		})
		if err != nil {
			return []survey.Question{}, []survey.SurveyResponse{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
		return qs, responses, nil
	} else {
		sql := `
		SELECT
			id,
			name,
			email,
			phone_number,
			rating,
			created_at
		FROM
			survey_respondents
		WHERE
			survey_id = $1
		ORDER BY
			created_at
		`
		rows, err := ss.db.Query(context.Background(), sql, s.Id)
		if err != nil {
			return []survey.Question{}, []survey.SurveyResponse{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
		defer rows.Close()

		responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (survey.SurveyResponse, error) {
			var response survey.SurveyResponse
			err := row.Scan(
				&response.Id,
				&response.Name,
				&response.Email,
				&response.PhoneNumber,
				&response.Rating,
				&response.CreatedAt,
			)
			response.Survey = s
			return response, err
		})
		if err != nil {
			return []survey.Question{}, []survey.SurveyResponse{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
		return qs, responses, nil
	}
}

func getTime(t time.Time) string {
	loc, _ := time.LoadLocation("America/Lima")
	locTime := t.In(loc)
	return locTime.Format(time.DateTime)
}

func (ss Survey) WriteResponses(w io.Writer, s survey.Survey) error {
	qs, sres, err := ss.GetResponses(s)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(w)

	// Cabecera
	fr := make([]string, 0, 5+len(qs))
	fr = append(fr, []string{
		"Nombre",
		"Correo",
		"Número",
		"Rating",
		"Respondido el",
	}...)
	for _, q := range qs {
		fr = append(fr, q.QuestionText)
	}
	if err := writer.Write(fr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Cuerpo
	for _, res := range sres {
		nr := make([]string, 0, 5+len(qs))
		nr = append(nr, []string{
			res.Name,
			res.Email,
			res.PhoneNumber,
			fmt.Sprintf("%d", res.Rating),
			getTime(res.CreatedAt),
		}...)
		nr = append(nr, res.Responses...)
		if err := writer.Write(nr); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	// Flush
	writer.Flush()
	if err := writer.Error(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func (ss Survey) InsertResponse(response survey.SurveyResponse) error {
	// Start transaction
	tx, err := ss.db.Begin(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer tx.Rollback(context.Background())

	sql := `
	INSERT INTO survey_respondents (
		survey_id,
		name,
		email,
		phone_number,
		rating
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)
	RETURNING
		id
	`

	// Insert respondent

	var respondentId int
	if err := tx.QueryRow(context.Background(), sql,
		response.Survey.Id,
		response.Name,
		response.Email,
		response.PhoneNumber,
		response.Rating,
	).Scan(&respondentId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Insert question responses

	for _, qResponse := range response.QuestionResponses {
		sql := `
		INSERT INTO question_responses (
			respondent_id,
			question_id,
			response_text
		)
		VALUES (
			$1,
			$2,
			$3
		)
		`
		_, err := tx.Exec(context.Background(), sql,
			respondentId,
			qResponse.Question.Id,
			qResponse.ResponseText,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// Image management

func (ss Survey) InsertImage(img *multipart.FileHeader) (store.Image, error) {
	// Source
	src, err := img.Open()
	if err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error opening the image")
	}
	defer src.Close()

	// Destination
	dst, err := os.CreateTemp(config.IMAGES_SAVEDIR, "*"+filepath.Ext(img.Filename))
	if err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error creating new image")
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error saving new image")
	}

	// Insert the image into database
	var newImg store.Image
	newImg.Filename = filepath.Base(dst.Name())
	if err := ss.db.QueryRow(context.Background(), `INSERT INTO images (filename)
VALUES ($1)
RETURNING id`, newImg.Filename).Scan(&newImg.Id); err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new image into database")
	}
	return newImg, nil
}

func (ss Survey) RemoveImage(id int) error {
	var filename string
	// Delete from database
	if err := ss.db.QueryRow(context.Background(), `DELETE FROM images WHERE id = $1 RETURNING filename`, id).Scan(&filename); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Image not found")
	}

	// Delete from filesystem
	if err := os.Remove(path.Join(config.IMAGES_SAVEDIR, filename)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error deleting the image from filesystem")
	}
	return nil
}

// Landing management

// TODO*
func (ss Survey) GetLandings() ([]survey.Landing, error) {
	sql := `
	SELECT
		l.id,
		l.title,
		l.content,
		l.created_at,
		l.is_published,
		s.id,
		s.title
	FROM
		landing l
		LEFT JOIN surveys s ON l.survey_id = s.id
	ORDER BY
		created_at DESC
	`
	rows, err := ss.db.Query(context.Background(), sql)
	if err != nil {
		return []survey.Landing{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()
	landings, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (survey.Landing, error) {
		var landing survey.Landing
		var sId *int
		var sTitle *string
		err := row.Scan(
			&landing.Id,
			&landing.Title,
			&landing.Content,
			&landing.CreatedAt,
			&landing.IsPublished,
			&sId,
			&sTitle,
		)
		if sId != nil {
			landing.Survey.Id = *sId
			landing.Survey.Title = *sTitle
		}
		return landing, err
	})
	if err != nil {
		return []survey.Landing{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return landings, nil
}

// TODO*
func (ss Survey) GetLandingById(id int) (landing survey.Landing, err error) {
	sql := `
	SELECT
		l.id,
		l.title,
		l.content,
		l.created_at,
		l.is_published,
		s.id,
		s.title
	FROM
		landing l
		LEFT JOIN surveys s ON l.survey_id = s.id
	WHERE
		l.id = $1
	`
	var sId *int
	var sTitle *string
	err = ss.db.QueryRow(context.Background(), sql, id).Scan(
		&landing.Id,
		&landing.Title,
		&landing.Content,
		&landing.CreatedAt,
		&landing.IsPublished,
		&sId,
		&sTitle,
	)
	if err != nil {
		err = echo.NewHTTPError(http.StatusNotFound, "Landing no encontrado")
		return
	}
	if sId != nil {
		landing.Survey.Id = *sId
		landing.Survey.Title = *sTitle
	}
	return
}

// TODO*
func (ss Survey) GetActiveLanding() (landing survey.Landing, err error) {
	sql := `
	SELECT
		l.id,
		l.title,
		l.content,
		l.created_at,
		l.is_published,
		s.id,
		s.title
	FROM
		landing l
		LEFT JOIN surveys s ON l.survey_id = s.id
	WHERE
		l.is_published = TRUE
	`
	var sId *int
	var sTitle *string
	err = ss.db.QueryRow(context.Background(), sql).Scan(
		&landing.Id,
		&landing.Title,
		&landing.Content,
		&landing.CreatedAt,
		&landing.IsPublished,
		&sId,
		&sTitle,
	)
	if err != nil {
		err = echo.NewHTTPError(http.StatusNotFound, "No hay un landing activo")
		return
	}
	if sId != nil {
		landing.Survey.Id = *sId
		landing.Survey.Title = *sTitle
	}
	return
}

func (ss Survey) InsertLanding(l survey.Landing) (id int, err error) {
	if l.Survey.Id != 0 {
		sql := `
		INSERT INTO landing (
			title,
			content,
			survey_id
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id
		`
		err = ss.db.QueryRow(context.Background(), sql, l.Title, l.Content, l.Survey.Id).Scan(&id)
	} else {
		sql := `
		INSERT INTO landing (
			title,
			content
		)
		VALUES (
			$1,
			$2
		)
		RETURNING
			id
		`
		err = ss.db.QueryRow(context.Background(), sql, l.Title, l.Content).Scan(&id)
	}
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	return
}

func (ss Survey) UpdateLanding(id int, uptLanding survey.Landing) (err error) {
	if uptLanding.Survey.Id != 0 {
		sql := `
		UPDATE landing
		SET
			title = $2,
			content = $3,
			survey_id = $4
		WHERE
			id = $1
		`
		_, err = ss.db.Exec(context.Background(), sql, id,
			uptLanding.Title,
			uptLanding.Content,
			uptLanding.Survey.Id,
		)
	} else {
		sql := `
		UPDATE landing
		SET
			title = $2,
			content = $3,
			survey_id = NULL
		WHERE
			id = $1
		`
		_, err = ss.db.Exec(context.Background(), sql, id,
			uptLanding.Title,
			uptLanding.Content,
		)
	}
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	return
}

func (ss Survey) RemoveLanding(l survey.Landing) (err error) {
	// Remove attached images
	imgs, err := ss.GetLandingImages(l)
	if err != nil {
		return
	}
	for _, img := range imgs {
		err := ss.RemoveImage(img.Id)
		if err != nil {
			continue
		}
	}
	// Remove landing from database
	sql := `DELETE FROM landing WHERE id = $1`
	_, err = ss.db.Exec(context.Background(), sql, l.Id)
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	return
}

func (ss Survey) SetActiveLanding(id int) (err error) {
	sql := `
	UPDATE landing
	SET
		is_published = FALSE
	WHERE
		is_published = TRUE
	`
	_, err = ss.db.Exec(context.Background(), sql)
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}

	sql1 := `
	UPDATE landing
	SET
		is_published = TRUE
	WHERE
		id = $1
	`
	_, err = ss.db.Exec(context.Background(), sql1, id)
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	return
}

func (ss Survey) HideLanding(id int) (err error) {
	sql := `
	UPDATE landing
	SET
		is_published = FALSE
	WHERE
		id = $1
	`
	_, err = ss.db.Exec(context.Background(), sql, id)
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	return
}

// Landing images management

func (ss Survey) GetLandingImages(l survey.Landing) (imgs []store.Image, err error) {
	sql := `
	SELECT
		img.id,
		img.filename,
		li.index
	FROM
		landing_images li
		JOIN images img ON li.image_id = img.id
	WHERE
		li.landing_id = $1
	ORDER BY
		li.index ASC, img.id ASC
	`
	rows, err := ss.db.Query(context.Background(), sql, l.Id)
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	imgs, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (img store.Image, err error) {
		err = row.Scan(&img.Id, &img.Filename, &img.Index)
		return
	})
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	return
}

func (ss Survey) GetMaxIndex(l survey.Landing) (int, error) {
	sql := `
	SELECT
		MAX(index)
	FROM
		landing_images
	WHERE
		landing_id = $1
	`
	var maxIndex *int
	err := ss.db.QueryRow(context.Background(), sql, l.Id).Scan(&maxIndex)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	if maxIndex == nil {
		return -1, nil
	}
	return *maxIndex, nil
}

func (ss Survey) ModifyLandingImages(l survey.Landing, imgs []store.Image) (err error) {
	// Start transaction
	tx, err := ss.db.Begin(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer tx.Rollback(context.Background())

	for _, img := range imgs {
		sql := `SELECT EXISTS (
			SELECT 1 FROM landing_images
			WHERE landing_id = $1 AND image_id = $2
		)`
		var exists bool
		if err := tx.QueryRow(context.Background(), sql, l.Id, img.Id).Scan(&exists); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if !exists {
			sql1 := `INSERT INTO landing_images (landing_id, image_id, index)
			VALUES ($1, $2, $3)`
			if _, err := tx.Exec(context.Background(), sql1, l.Id, img.Id, img.Index); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		} else {
			sql2 := `UPDATE landing_images
			SET index = $1
			WHERE landing_id = $2 AND image_id = $3`
			if _, err := tx.Exec(context.Background(), sql2, img.Index, l.Id, img.Id); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
