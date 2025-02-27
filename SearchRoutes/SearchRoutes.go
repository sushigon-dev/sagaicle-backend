package SearchRoutes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type sort struct {
	KEY   string `json: "KEY"`
	ORDER string `json: "ORDER"`
}

type request struct {
	DISTANCE      float64  `json: "DISTANCE"`
	TIME          int64    `json: "TIME"`
	TAGS          []string `json: "TAGS"`
	SEARCH_OPTION []string `json: "SEARCH_OPTION"`
	SORT          sort     `json: "SORT"`
	LIMIT         int64    `json: "LIMIT"`
}

type response struct {
	ID          int64    `json: "ID"`
	TITLE       string   `json: "TITLE"`
	DESCRIPTION string   `json: "DESCRIPTION"`
	DISTANCE    float64  `json: "DISTANCE"`
	TIME        int64    `json: "TIME"`
	TAGS        []string `json: "TAGS"`
	LIKES       int64    `json: "LIKES"`
	IMAGE       string   `json: "IMAGE"`
	UPDATE_AT   string   `json: "UPDATE_AT"`
}

type TagHandler struct {
	DB *sql.DB
}

func (h *SearchHandler) SearchRoutes(c *gin.Context) {

}
