package display

import (
	"log"
	"net/http"
	queries "sharing_vision/internal/db"
	"sharing_vision/pages"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Display struct {
	query queries.Queries
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func NewDisplayController(query queries.Queries) *Display {
	return &Display{
		query: query,
	}
}

func (d *Display) AllPostDisplay(c echo.Context) error {
	limit := c.Param("limit")
	offset := c.Param("offset")
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error parsing query param to integer, error :")
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Somethings wrong in the server. please try again later."})
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Println("Error parsing query param to integer, error :")
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Somethings wrong in the server. please try again later."})
	}

	filterQueries := queries.ListArticlesParams{
		Limit:  int32(limitInt),
		Offset: int32(offsetInt),
	}

	res, err := d.query.ListArticles(c.Request().Context(), filterQueries)
	if err != nil {
		log.Println("Error fetching articles with filters,error :")
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to fetch request."})
	}
	allPost := pages.AllPosts(res)
	return Render(c, http.StatusOK, pages.Index(allPost))
}
