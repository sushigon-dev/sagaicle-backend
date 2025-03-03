package SearchRoutes

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type SearchRequest struct {
	Distance struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"distance"`

	Time struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"time"`

	Tags         []string `json:"tags"`
	SearchOption string   `json:"search_option"`
	Sort         struct {
		Key   string `json:"key"`
		Order string `json:"order"`
	} `json:"sort"`
	Limit int `json:"limit"`
}

// ルート情報のレスポンス用構造体
type PreRoute struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Full_description  string  `json:"full_description"`
	Distance          float64 `json:"distance"`
	Time              int     `json:"time"`
	Tags              string  `json:"tags"`
	Total_Checkpoints int     `json:"total_checkpoints"`
	Images            string  `json:"images"`
	Map               string  `json:"map"`
	Name              string  `json:"name"`
	Lat               int     `json:"lat"`
	Ing               int     `json:"ing"`
	Route_ID          string  `json:"route_id"`
	UpdateAt          string  `json:"update_at"`
	Image             string  `json:"image"`
	Likes             int     `json:"likes"`
}
type Route struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Distance    float64 `json:"distance"`
	Time        int     `json:"time"`
	Tags        string  `json:"tags"`
	Likes       int     `json:"likes"`
	Image       string  `json:"image"`
	UpdateAt    string  `json:"update_at"`
}

// レスポンスの構造体
type SearchResponse struct {
	HitCount int           `json:"hit_count"`
	Routes   []Route       `json:"routes"`
	Request  SearchRequest `json:"request"`
}

type SearchHandler struct {
	DB *sql.DB
}

func (h *SearchHandler) SearchRoutes(c *gin.Context) {
	// リクエストデータをバインド
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	switch req.SearchOption {
	case "AND":
		query := "SELECT * FROM routes WHERE (DISTANCE BETWEEN ? AND ?) AND (TIME BETWEEN ? AND ?)"
		args := []interface{}{req.Distance.Min, req.Distance.Max, req.Time.Min, req.Time.Max}

		if len(req.Tags) > 0 {
			tagConditions := []string{}
			for _, tag := range req.Tags {
				tagConditions = append(tagConditions, "tags LIKE ?")
				args = append(args, "%"+tag+"%")
			}
			query += "AND (" + strings.Join(tagConditions, " OR ") + ")"
		}

		// 取得制限
		limit := 12
		if req.Limit >= 1 && req.Limit <= 60 {
			limit = req.Limit
		}
		query += " LIMIT ?"
		args = append(args, limit)

		rows, err := h.DB.Query(query, args...)
		if err != nil {
			log.Println("Query Error:", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}
		defer rows.Close()

		var routes []Route
		for rows.Next() {
			var route Route
			var preroute PreRoute
			if err := rows.Scan(&preroute.ID, &preroute.Title, &preroute.Description, &preroute.Full_description, &preroute.Distance,
				&preroute.Time, &preroute.Tags, &preroute.Total_Checkpoints, &preroute.Images, &preroute.Map, &preroute.Name, &preroute.Lat,
				&preroute.Ing, &preroute.Route_ID, &preroute.UpdateAt, &preroute.Image, &preroute.Likes); err != nil {
				log.Println("Scan Error:", err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Data scan error"})
				return
			}
			route.ID = preroute.ID
			route.Title = preroute.Title
			route.Description = preroute.Description
			route.Distance = preroute.Distance
			route.Time = preroute.Time
			route.Tags = preroute.Tags
			route.Likes = preroute.Likes
			route.Image = preroute.Image
			route.UpdateAt = preroute.UpdateAt

			routes = append(routes, route)
		}
		defer rows.Close()

		// レスポンスを返す
		response := SearchResponse{
			HitCount: len(routes),
			Routes:   routes,
			Request:  req,
		}
		c.IndentedJSON(http.StatusOK, response)

	case "OR":
		query := "SELECT * FROM routes WHERE (DISTANCE BETWEEN ? AND ?) OR (TIME BETWEEN ? AND ?)"
		args := []interface{}{req.Distance.Min, req.Distance.Max, req.Time.Min, req.Time.Max}

		if len(req.Tags) > 0 {
			tagConditions := []string{}
			for _, tag := range req.Tags {
				tagConditions = append(tagConditions, "tags LIKE ?")
				args = append(args, "%"+tag+"%")
			}
			query += "OR (" + strings.Join(tagConditions, " OR ") + ")"
		}

		rows, err := h.DB.Query(query, args...)
		if err != nil {
			log.Println("Query Error:", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}
		defer rows.Close()

		var routes []Route
		for rows.Next() {
			var route Route
			var preroute PreRoute
			if err := rows.Scan(&preroute.ID, &preroute.Title, &preroute.Description, &preroute.Full_description, &preroute.Distance,
				&preroute.Time, &preroute.Tags, &preroute.Total_Checkpoints, &preroute.Images, &preroute.Map, &preroute.Name, &preroute.Lat,
				&preroute.Ing, &preroute.Route_ID, &preroute.UpdateAt, &preroute.Image, &preroute.Likes); err != nil {
				log.Println("Scan Error:", err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Data scan error" + err.Error()})
				return
			}
			route.ID = preroute.ID
			route.Title = preroute.Title
			route.Description = preroute.Description
			route.Distance = preroute.Distance
			route.Time = preroute.Time
			route.Tags = preroute.Tags
			route.Likes = preroute.Likes
			route.Image = preroute.Image
			route.UpdateAt = preroute.UpdateAt

			routes = append(routes, route)
		}

		// レスポンスを返す
		response := SearchResponse{
			HitCount: len(routes),
			Routes:   routes,
			Request:  req,
		}
		c.IndentedJSON(http.StatusOK, response)

	}

}

/* リクエスト例 */
/*
curl -X POST "http://localhost:8080/api/search" \
     -H "Content-Type: application/json" \
     -d '{
           "distance": { "min": 0, "max": 20000},
           "time": { "min": 0, "max": 1800 },
           "tags": ["mountain", "scenic", "hiking", "nature"],
           "search_option": "OR",
           "sort": { "order": "desc" },
           "limit": 10
         }'
*/
