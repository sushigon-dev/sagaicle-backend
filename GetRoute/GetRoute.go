package GetRoute

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type RouteHandler struct {
	DB *sql.DB
}

type route struct {
	ID               string `json:"id"`
	TITLE            string `json:"title"`
	DESCRIPTION      string `json:"description"`
	FULL_DESCRIPTION string `json:"full_description"`
	DISTANCE         int64  `json:"distance"`
	TIME             int64  `json:"time"`
	TAGS             string `json:"tags"`
	CHECKPOINT_COUNT int64  `json:"checkpoint_count"`
	IMAGES           string `json:"images"`
	MAP              string `json:"map"`
	UPDATE           string `json:"update"`
}

func (h *RouteHandler) GetRoute(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	var receive route

	err := h.DB.QueryRow("SELECT ID, TITLE, DESCRIPTION, FULL_DESCRIPTION, DISTANCE, TIME, TAGS, TOTAL_CHECKPOINTS, IMAGES, MAP, UPDATE_AT FROM routes WHERE ID = ?", id).
		Scan(&receive.ID, &receive.TITLE, &receive.DESCRIPTION, &receive.FULL_DESCRIPTION,
			&receive.DISTANCE, &receive.TIME, &receive.TAGS, &receive.CHECKPOINT_COUNT, &receive.IMAGES, &receive.MAP, &receive.UPDATE)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "route not found " + err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, receive)
}
