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
	ID                string  `json:"id"`
	TITLE             string  `json:"title"`
	DESCRIPTION       string  `json:"description"`
	FULL_DESCRIPTION  string  `json:"full_description"`
	DISTANCE          float64 `json:"distance"`
	TIME              int64   `json:"time"`
	TAGS              string  `json:"tags"`
	TOTAL_CHECKPOINTS int64   `json:"checkpoint_count"`
	IMAGES            string  `json:"images"`
	MAP               string  `json:"map"`
	UPDATE_AT         string  `json:"update_at"`
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
			&receive.DISTANCE, &receive.TIME, &receive.TAGS, &receive.TOTAL_CHECKPOINTS, &receive.IMAGES, &receive.MAP, &receive.UPDATE_AT)

	if err != nil {
		//エラーを表示すること
		c.JSON(http.StatusInternalServerError, gin.H{"error": "route not found" + err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, receive)
}

//GetRouteへのリクエスト時は
//idはシングルクォ-トでかこむ
//例：SELECT ID, TITLE, DESCRIPTION, FULL_DESCRIPTION, DISTANCE, TIME, TAGS, TOTAL_CHECKPOINTS, IMAGES, MAP, UPDATE_AT FROM routes WHERE ID = '0a929a67-d964-4c1b-aa8a-12ddc08ed957';
