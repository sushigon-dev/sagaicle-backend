package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sushigon-dev/sagaicle-backend/GetRoute"
	"github.com/sushigon-dev/sagaicle-backend/GetTags"
	//"github.com/sushigon-dev/sagaicle-backend/GetRoute"
)

func main() {

	db, err := sql.Open("sqlite3", "./GetRoute.db")

	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	tagHandler := &GetTags.TagHandler{DB: db}
	routeHandler := &GetRoute.RouteHandler{DB: db}

	router := gin.Default()
	router.GET("/api/tags", tagHandler.GetTags)
	router.GET("/api/route/:id", routeHandler.GetRoute)

	router.Run("localhost:8080")
}
