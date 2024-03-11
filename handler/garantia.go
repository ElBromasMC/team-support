package handler

import (
	"alc/model"
	"alc/view/garantia"
	"net/http"

	"github.com/labstack/echo/v4"
)

var garantiaSubCategories []model.StoreSubCategory
var garantiaItems []model.StoreItem

func init() {
	// garantiaSubCategories = []model.StoreSubCategory{
	// 	{Name: "NB (Escritorio)", Description: "ZenBook - VivoBook - AsusLaptop", Img: "/static/img/NB.png", Slug: "nb"},
	// 	{Name: "NR (Gamer)", Description: "TUF - ROG", Img: "/static/img/NR.png", Slug: "nr"},
	// 	{Name: "PT (All in one)", Description: "All in one", Img: "/static/img/PT.png", Slug: "pt"},
	// }

	// newUuid0, _ := uuid.NewV4()
	// newUuid1, _ := uuid.NewV4()
	// newUuid2, _ := uuid.NewV4()
	// item0 := model.StoreItem{
	// 	Uuid:        newUuid0,
	// 	Category:    "GARANTIA",
	// 	SubCategory: "NB (Escritorio)",
	// 	Name:        "Garantía + Daño Accidental + Bateria TUF",
	// 	Price:       18000,
	// 	Slug:        "garantia-accidental-bateria",
	// 	Img:         "/static/img/garantia1.jpg",
	// }
	// item1 := model.StoreItem{
	// 	Uuid:        newUuid1,
	// 	Category:    "GARANTIA",
	// 	SubCategory: "NR (Gamer)",
	// 	Name:        "Protección contra daño accidental TUF",
	// 	Price:       18000,
	// 	Slug:        "accidental",
	// 	Img:         "/static/img/garantia1.jpg",
	// }
	// item2 := model.StoreItem{
	// 	Uuid:        newUuid2,
	// 	Category:    "GARANTIA",
	// 	SubCategory: "PT (All in one)",
	// 	Name:        "Garantía Extendida + Domicilio TUF",
	// 	Price:       18000,
	// 	Slug:        "garantia-domicilio",
	// 	Img:         "/static/img/garantia1.jpg",
	// }
	// garantiaItems = append(garantiaItems, item0, item0, item0, item0, item0, item0, item0, item0, item0, item1, item2)
}

func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	return render(c, http.StatusOK, garantia.Show(garantiaSubCategories))
}

func (h *Handler) HandleGarantiaCategoryShow(c echo.Context) error {
	// Category slug from path `/garantia/:slug`
	slug := c.Param("slug")

	found := false
	var category model.StoreSubCategory
	for _, i := range garantiaSubCategories {
		if i.Slug == slug {
			category = i
			found = true
			break
		}
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	items := []model.StoreItem{}
	for _, i := range garantiaItems {
		if i.SubCategory == category.Name {
			items = append(items, i)
		}
	}

	return render(c, http.StatusOK, garantia.ShowCategory(category, items))
}

func (h *Handler) HandleGarantiaItemShow(c echo.Context) error {
	// Item slug from path `/garantia/:categorySlug/:itemSlug`
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	var found bool

	found = false
	var category model.StoreSubCategory
	for _, i := range garantiaSubCategories {
		if i.Slug == categorySlug {
			category = i
			found = true
			break
		}
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	found = false
	var item model.StoreItem
	for _, i := range garantiaItems {
		if i.SubCategory == category.Name && i.Slug == itemSlug {
			item = i
			found = true
			break
		}
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return render(c, http.StatusOK, garantia.ShowItem(item))
}
