package GetTags

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func GetTags(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tags []string

		rows, err := db.Query("SELECT TAG FROM tags")
		if err != nil {
			log.Println("Query Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				log.Println("Scan Error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Data scan error"})
				return
			}
			tags = append(tags, tag)
		}
		c.JSON(http.StatusOK, gin.H{"tags": tags})
	}
}
