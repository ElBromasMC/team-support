package handler

import (
	"alc/config"
	"alc/model/store"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, statusCode int, t templ.Component) error {
	c.Response().Writer.WriteHeader(statusCode)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response().Writer)
}

func RemoveCookie(c echo.Context, name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.MaxAge = -1
	c.SetCookie(cookie)
}

func (h *Handler) InsertImage(img *multipart.FileHeader) (int, error) {
	// Source
	src, err := img.Open()
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error opening the image")
	}
	defer src.Close()

	// Destination
	dst, err := os.CreateTemp(config.IMAGES_SAVEDIR, "*"+filepath.Ext(img.Filename))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error creating new image")
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error saving new image")
	}

	var imgId int
	// Insert the image in the database
	if err := h.DB.QueryRow(context.Background(), `INSERT INTO images (filename)
VALUES ($1)
RETURNING id`, filepath.Base(dst.Name())).Scan(&imgId); err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new image into database")
	}
	return imgId, nil
}

func (h *Handler) GetCategories(t store.Type) ([]store.Category, error) {
	rows, err := h.DB.Query(context.Background(), `SELECT sc.id, sc.name, sc.description, sc.slug, img.filename
FROM store_categories AS sc
LEFT JOIN images AS img
ON sc.img_id = img.id
WHERE type = $1`, t)
	if err != nil {
		return []store.Category{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	cats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.Category, error) {
		var cat store.Category
		var img *string
		err := row.Scan(&cat.Id, &cat.Name, &cat.Description, &cat.Slug, &img)
		if img != nil {
			cat.Img.Filename = *img
		} else {
			cat.Img.Filename = ""
		}
		cat.Type = store.GarantiaType
		return cat, err
	})
	if err != nil {
		return []store.Category{}, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return cats, nil
}
