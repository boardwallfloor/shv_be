package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	queries "sharing_vision/internal/db"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type Server struct {
	db queries.Queries
}

func NewServer(queries queries.Queries) *Server {
	return &Server{
		db: queries,
	}
}

type CreateArticleRequest struct {
	Title    string  `json:"title"`
	Content  *string `json:"content"`
	Category *string `json:"category"`
	Status   string  `json:"status"`
}

func (s *Server) AddArticle(c echo.Context) error {
	log.Println("Adding article")
	req := new(CreateArticleRequest)
	if err := c.Bind(req); err != nil {
		log.Println("Error parsing request body, error :", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Unable to parse request input."})
	}

	if len(req.Title) < 20 {
		log.Println("Validation error, title is too short:", len(req.Title))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Title is too short, min 20 character"})
	}
	if req.Content == nil || len(*req.Content) < 200 {
		log.Println("Validation error, content is missing or too short")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Content is required and must be at least 200 characters"})
	}
	if req.Category == nil || len(*req.Category) < 3 {
		log.Println("Validation error, category is missing or too short")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Category is required and must be at least 3 characters"})
	}

	params := queries.CreateArticleParams{
		Title:  req.Title,
		Status: queries.ArticlesStatus(req.Status),
	}
	if req.Content != nil {
		params.Content = null.StringFrom(*req.Content)
	}
	if req.Category != nil {
		params.Category = null.StringFrom(*req.Category)
	}

	_, err := s.db.CreateArticle(c.Request().Context(), params)
	if err != nil {
		log.Println("Error inserting new articles, error :", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to save new articles."})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Article created successfully"})
}

func (s *Server) GetArticles(c echo.Context) error {
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

	res, err := s.db.ListArticles(c.Request().Context(), filterQueries)
	if err != nil {
		log.Println("Error fetching articles with filters,error :")
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to fetch request."})
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "Success", "data": res})
}

func (s *Server) GetArticleById(c echo.Context) error {
	articleId := c.Param("id")
	articleIdInt, err := strconv.Atoi(articleId)
	if err != nil {
		log.Println("Error parsing article id")
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID parameter"})
	}

	res, err := s.db.GetArticle(c.Request().Context(), int32(articleIdInt))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows with matching id")
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Unable to find articles with id"})
		} else {
			log.Println("Error fetching request with id")
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to fetch request."})
		}
	}
	return c.JSON(http.StatusOK, map[string]any{"message": "Success", "data": res})
}

type UpdateArticleRequest struct {
	Title    string  `json:"title"`
	Content  *string `json:"content"`
	Category *string `json:"category"`
	Status   string  `json:"status"`
}

func (s *Server) UpdateArticleById(c echo.Context) error {
	articleId := c.Param("id")
	articleIdInt, err := strconv.Atoi(articleId)
	if err != nil {
		log.Println("Error parsing article id:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID parameter"})
	}

	req := new(UpdateArticleRequest)
	if err := c.Bind(req); err != nil {
		log.Println("Error parsing request body, error :", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Unable to parse request input."})
	}

	if len(req.Title) < 20 {
		log.Println("Validation error, title is too short")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Title is too short, min 20 character"})
	}

	if req.Content == nil || len(*req.Content) < 200 {
		log.Println("Validation error, content is missing or too short")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Content is required and must be at least 200 characters"})
	}

	if req.Category == nil || len(*req.Category) < 3 {
		log.Println("Validation error, category is missing or too short")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Category is required and must be at least 3 characters"})
	}

	params := queries.UpdateArticleParams{
		ID:     int32(articleIdInt),
		Title:  req.Title,
		Status: queries.ArticlesStatus(req.Status),
	}

	if req.Content != nil {
		params.Content = null.StringFrom(*req.Content)
	}
	if req.Category != nil {
		params.Category = null.StringFrom(*req.Category)
	}

	err = s.db.UpdateArticle(c.Request().Context(), params)
	if err != nil {
		log.Println("Error updating articles, error :", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to update articles"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Article updated successfully"})
}

func (s *Server) DeleteArticleById(c echo.Context) error {
	articleId := c.Param("id")
	articleIdInt, err := strconv.Atoi(articleId)
	if err != nil {
		log.Println("Error parsing article id")
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID parameter"})
	}

	err = s.db.DeleteArticle(c.Request().Context(), int32(articleIdInt))
	if err != nil {
		log.Println("Error deleting article with id, error :")
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to delete article"})
	}
	return c.JSON(http.StatusOK, map[string]string{})
}
