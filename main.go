package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sushigon-dev/sagaicle-backend/GetRoute"
	"github.com/sushigon-dev/sagaicle-backend/GetTags"
	"github.com/sushigon-dev/sagaicle-backend/PostRoute"
	"github.com/sushigon-dev/sagaicle-backend/SearchRoutes"
)

func main() {

	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	tagHandler := &GetTags.TagHandler{DB: db}
	routeHandler := &GetRoute.RouteHandler{DB: db}
	searchHandler := &SearchRoutes.SearchHandler{DB: db}
	postHandler := &PostRoute.PostHandler{DB: db}

	router := gin.Default()
	router.GET("/api/tags", tagHandler.GetTags)
	router.GET("/api/route/:id", routeHandler.GetRoute)
	router.POST("/api/search", searchHandler.SearchRoutes)
	router.POST("/api/routes", postHandler.PostRoute)

	router.Run("localhost:8080")
}
