package main

import (
	"database/sql"
	"log"

	"sagaicle-backend/GetTags"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/api/tags", GetTags.GetTags(db))

	router.Run("localhost:8080")
}
