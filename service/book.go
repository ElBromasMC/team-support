package service

import (
	"alc/model/book"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/labstack/echo/v4"
)

type Book struct {
	db *pgxpool.Pool
}

func NewBookService(db *pgxpool.Pool) Book {
	return Book{
		db: db,
	}
}

func (bs Book) GetDocumentType(slug string) (book.DocumentType, error) {
	if slug == "DNI" {
		return book.DNIType, nil
	} else if slug == "CARNET" {
		return book.CarnetType, nil
	} else if slug == "OTHER" {
		return book.OtherType, nil
	}
	return "", echo.NewHTTPError(http.StatusBadRequest, "Tipo de documento inválido")
}

func (bs Book) GetGoodType(slug string) (book.GoodType, error) {
	if slug == "PRODUCT" {
		return book.ProductType, nil
	} else if slug == "SERVICE" {
		return book.ServiceType, nil
	}
	return "", echo.NewHTTPError(http.StatusBadRequest, "Tipo de bien inválido")
}

func (bs Book) GetComplaintType(slug string) (book.ComplaintType, error) {
	if slug == "RECLAMO" {
		return book.ReclamoType, nil
	} else if slug == "QUEJA" {
		return book.QuejaType, nil
	}
	return "", echo.NewHTTPError(http.StatusBadRequest, "Tipo de reclamo inválido")
}

func (bs Book) GetBookEntryById(id int) (entry book.Entry, err error) {
	sql := `
	SELECT
		id,
		created_at,
		document_type,
		document_number,
		name,
		address,
		phone_number,
		email,
		parent_name,
		good_type,
		good_description,
		complaint_type,
		complaint_description,
		actions_description
	FROM
		book_entries
	WHERE
		id = $1
	`
	err = bs.db.QueryRow(context.Background(), sql, id).Scan(
		&entry.Id,
		&entry.CreatedAt,
		&entry.DocumentType,
		&entry.DocumentNumber,
		&entry.Name,
		&entry.Address,
		&entry.PhoneNumber,
		&entry.Email,
		&entry.ParentName,
		&entry.GoodType,
		&entry.GoodDescription,
		&entry.ComplaintType,
		&entry.ComplaintDescription,
		&entry.ActionsDescription,
	)
	if err != nil {
		err = echo.NewHTTPError(http.StatusNotFound, "Entrada de libro no encontrada")
		return
	}
	return
}

func (bs Book) InsertBookEntry(entry book.Entry) (id int, err error) {
	sql := `
	INSERT INTO book_entries (
		document_type,
		document_number,
		name,
		address,
		phone_number,
		email,
		parent_name,
		good_type,
		good_description,
		complaint_type,
		complaint_description,
		actions_description
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12
	)
	RETURNING
		id
	`
	err = bs.db.QueryRow(context.Background(), sql,
		entry.DocumentType,
		entry.DocumentNumber,
		entry.Name,
		entry.Address,
		entry.PhoneNumber,
		entry.Email,
		entry.ParentName,
		entry.GoodType,
		entry.GoodDescription,
		entry.ComplaintType,
		entry.ComplaintDescription,
		entry.ActionsDescription,
	).Scan(&id)
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError)
		return
	}
	return
}

// PDF generation

func getDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getPageHeader() core.Row {
	headerFontSize := 8.0
	return row.New(20).Add(
		image.NewFromFileCol(3, "./assets/static/img/logo1.png", props.Rect{
			Center:  true,
			Percent: 80,
		}),
		col.New(3),
		col.New(6).Add(
			text.New("TEAM SUPPORT SERVICES S.A.C", props.Text{
				Size:  headerFontSize,
				Align: align.Right,
			}),
			text.New("RUC 20600355831", props.Text{
				Top:   4,
				Size:  headerFontSize,
				Align: align.Right,
			}),
			text.New("AV. MONTE DE LOS OLIVOS 993, SANTIAGO DE SURCO - LIMA", props.Text{
				Top:   12,
				Size:  headerFontSize,
				Align: align.Right,
			}),
		),
	)
}

func getDateOnly(t time.Time) string {
	loc, _ := time.LoadLocation("America/Lima")
	locTime := t.In(loc)
	return locTime.Format("02/01/2006")
}

func (bs Book) GeneratePDF(entry book.Entry) (core.Maroto, error) {
	grayColor := getGrayColor()
	darkGrayColor := getDarkGrayColor()
	bodyTitleFontSize := 10.0
	bodyFontSize := 9.0
	leftSeparation := 1.5
	bodyProps := props.Text{
		Top:  1.5,
		Left: leftSeparation,
		Size: bodyFontSize,
	}

	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	// Header
	err := m.RegisterHeader(getPageHeader())
	if err != nil {
		return m, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Body title
	m.AddRows(text.NewRow(10, "TEAM SUPPORT SERVICES: LIBRO DE RECLAMACIONES", props.Text{
		Top:   3,
		Size:  bodyTitleFontSize,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	m.AddRow(8,
		text.NewCol(12, "HOJA DE RECLAMACIÓN", props.Text{
			Top:   1.5,
			Size:  bodyTitleFontSize,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGrayColor})
	m.AddRow(8,
		text.NewCol(12, fmt.Sprintf("FECHA: %s - NRO: %d-%d", getDateOnly(entry.CreatedAt), entry.Id, entry.CreatedAt.Year()), props.Text{
			Top:   1.5,
			Size:  bodyTitleFontSize,
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	// Body
	m.AddRow(16,
		text.NewCol(12, "1. IDENTIFICACIÓN DEL CONSUMIDOR RECLAMANTE", props.Text{
			Top:   9.5,
			Left:  leftSeparation,
			Size:  bodyFontSize,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGrayColor})
	m.AddRow(8,
		text.NewCol(3, "NOMBRE:", bodyProps),
		text.NewCol(4, entry.Name, bodyProps),
		text.NewCol(1, entry.DocumentType.ToSlug()+":", bodyProps),
		text.NewCol(4, entry.DocumentNumber, bodyProps),
	)
	m.AddRow(8,
		text.NewCol(3, "DOMICILIO:", bodyProps),
		text.NewCol(9, entry.Address, bodyProps),
	).WithStyle(&props.Cell{BackgroundColor: grayColor})
	m.AddRow(8,
		text.NewCol(3, "TELÉFONO:", bodyProps),
		text.NewCol(4, entry.PhoneNumber, bodyProps),
		text.NewCol(1, "E-MAIL:", bodyProps),
		text.NewCol(4, entry.Email, bodyProps),
	)
	m.AddRow(8,
		text.NewCol(3, "PADRE O MADRE:", bodyProps),
		text.NewCol(9, entry.ParentName, bodyProps),
	).WithStyle(&props.Cell{BackgroundColor: grayColor})
	m.AddRow(16,
		text.NewCol(12, "2. IDENTIFICACIÓN DEL BIEN CONTRATADO", props.Text{
			Top:   9.5,
			Left:  leftSeparation,
			Size:  bodyFontSize,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGrayColor})
	m.AddRow(8,
		text.NewCol(3, "TIPO BIEN CONTRATADO:", bodyProps),
		text.NewCol(9, strings.ToUpper(entry.GoodType.ToSlug()), bodyProps),
	)
	m.AddRow(8,
		text.NewCol(3, "DESCRIPCIÓN:", bodyProps),
		text.NewCol(9, entry.GoodDescription, bodyProps),
	).WithStyle(&props.Cell{BackgroundColor: grayColor})
	m.AddRow(16,
		text.NewCol(12, "3. DETALLE DE LA RECLAMACIÓN", props.Text{
			Top:   9.5,
			Left:  leftSeparation,
			Size:  bodyFontSize,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGrayColor})
	m.AddRow(8,
		text.NewCol(3, "TIPO DE LA RECLAMACIÓN:", bodyProps),
		text.NewCol(9, strings.ToUpper(entry.ComplaintType.ToSlug()), bodyProps),
	)
	m.AddRow(24,
		text.NewCol(3, "DETALLE:", props.Text{
			Top:  10,
			Left: leftSeparation,
			Size: bodyFontSize,
		}),
		text.NewCol(9, entry.ComplaintDescription, bodyProps),
	).WithStyle(&props.Cell{BackgroundColor: grayColor})
	m.AddRow(16,
		text.NewCol(12, "4. ACCIONES ADOPTADAS POR EL PROVEEDOR", props.Text{
			Top:   9.5,
			Left:  leftSeparation,
			Size:  bodyFontSize,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGrayColor})
	m.AddRow(8,
		text.NewCol(3, "DETALLE:", bodyProps),
		text.NewCol(9, entry.ActionsDescription, bodyProps),
	)
	return m, nil
}
