package PostRoute

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ルートデータの構造体
type Route struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	FullDescription  string       `json:"full_description"`
	Distance         float64      `json:"distance"`
	Time             int64        `json:"time"`
	Tags             []string     `json:"tags"`
	TotalCheckpoints int          `json:"total_checkpoints"`
	Images           []string     `json:"images"`
	Map              string       `json:"map"`
	Checkpoints      []Checkpoint `json:"checkpoints"`
}

// チェックポイントの構造体
type Checkpoint struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

type PostHandler struct {
	DB *sql.DB
}

// `POST /api/routes` エンドポイント
func (h *PostHandler) PostRoute(c *gin.Context) {
	var route Route

	// JSON を構造体にバインド
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// UUIDを生成
	route.ID = uuid.New().String()

	// データベースに格納
	_, err := h.DB.Exec(`
		INSERT INTO routes (ID, TITLE, DESCRIPTION, FULL_DESCRIPTION, DISTANCE, TIME, TOTAL_CHECKPOINTS, MAP)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		route.ID, route.Title, route.Description, route.FullDescription, route.Distance, route.Time, route.TotalCheckpoints, route.Map,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert route"})
		return
	}

	// タグをデータベースに保存
	for _, tag := range route.Tags {
		_, err := h.DB.Exec(`INSERT INTO routes (ROUTE_ID, TAG) VALUES (?, ?)`, route.ID, tag)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert tags"})
			return
		}
	}

	// 画像URLをデータベースに保存
	for _, img := range route.Images {
		_, err := h.DB.Exec(`INSERT INTO routes (ROUTE_ID, IMAGE_URL) VALUES (?, ?)`, route.ID, img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert images"})
			return
		}
	}

	// チェックポイントをデータベースに保存
	for _, cp := range route.Checkpoints {
		_, err := h.DB.Exec(`INSERT INTO routes (ROUTE_ID, NAME, LAT, LNG) VALUES (?, ?, ?, ?)`, route.ID, cp.Name, cp.Lat, cp.Lng)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert checkpoints"})
			return
		}
	}

	// 受け取ったデータをそのままレスポンスとして返す
	c.JSON(http.StatusOK, route)
}
