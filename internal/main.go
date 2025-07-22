package main

import (
	"database/sql"
	"log"
	"sharing_vision/internal/api"
	queries "sharing_vision/internal/db"
	"sharing_vision/internal/display"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetDBConn() *queries.Queries {
	db, err := sql.Open("mysql", "root:very_sekrit@/shv_db?parseTime=true")
	if err != nil {
		log.Println("Unable to open connection to MySQL DB")
		log.Fatal(err)
	}
	query := queries.New(db)
	return query
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	app := echo.New()
	dbConn := SetDBConn()

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
	}))

	displayController := display.NewDisplayController(*dbConn)
	displayRoutes := app.Group("/display")
	displayRoutes.GET("/", displayController.AllPostDisplay)

	apiServer := api.NewServer(*dbConn)
	apiRoutes := app.Group("/article")
	apiRoutes.POST("/", apiServer.AddArticle)
	apiRoutes.GET("/:id", apiServer.GetArticleById)
	apiRoutes.GET("/:limit/:offest", apiServer.GetArticles)
	apiRoutes.PUT("/:id", apiServer.UpdateArticleById)
	apiRoutes.DELETE("/:id", apiServer.DeleteArticleById)

	app.Logger.Fatal(app.Start(":8080"))
}
