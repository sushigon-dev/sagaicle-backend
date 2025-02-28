package PostRoute

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ルートデータの構造体
type Route struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	FullDescription  string     `json:"full_description"`
	Distance         float64    `json:"distance"`
	Time             int64      `json:"time"`
	Tags             []string   `json:"tags"`
	TotalCheckpoints int        `json:"total_checkpoints"`
	Images           []string   `json:"images"`
	Map              string     `json:"map"`
	Checkpoints      []struct { // ✅ 配列（スライス）に変更
		Name string  `json:"name"`
		Lat  float64 `json:"lat"`
		Ing  float64 `json:"lng"`
	} `json:"checkpoints"`
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
	//　ただしGetRouteなどからアクセスしようとすると
	// 保存されていない、と出る。
	tagsString := strings.Join(route.Tags, ",")
	_, err = h.DB.Exec(`INSERT INTO routes (ROUTE_ID, TAGS) VALUES (?, ?)`, route.ID, tagsString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert tags" + err.Error()})
		return
	}

	// 画像URLをデータベースに保存
	imgString := strings.Join(route.Images, ",")
	_, err = h.DB.Exec(`INSERT OR REPLACE INTO routes (ID, IMAGES) VALUES (?, ?)`, route.ID, imgString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert tags" + err.Error()})
		return
	}

	// チェックポイントをデータベースに保存
	for _, cp := range route.Checkpoints {
		_, err := h.DB.Exec(
			`INSERT INTO routes (ROUTE_ID, NAME, LAT, ING) VALUES (?, ?, ?, ?)`,
			route.ID, cp.Name, cp.Lat, cp.Ing,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert checkpoints"})
			return
		}
	}

	/*
		_, err = h.DB.Exec(`INSERT INTO routes (ROUTE_ID, NAME, LAT, ING) VALUES (?, ?, ?, ?)`, route.ID, route.Checkpoints.Name, route.Checkpoints.Lat, route.Checkpoints.Ing)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert checkpoints"})
			return
		}
	*/

	// 受け取ったデータをそのままレスポンスとして返す
	c.JSON(http.StatusOK, route)
}

//リクエスト例
/*
curl -X POST "http://localhost:8080/api/routes" \
-H "Authorization: Bearer <YOUR_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "title": "サンプルルート",
  "description": "楽しいサイクリングコース",
  "full_description": "このルートは海岸沿いを走る爽快なサイクリングコースです。",
  "distance": 12.5,
  "time": 90,
  "tags": ["海", "山", "初心者向け"],
  "total_checkpoints": 3,
  "images": [
	"https://github.com/sushigon-dev/sagaicle-docs/images/0e746155-f52e-4502-83f8-ab91e47abf6f.png",
	"https://x162-43-27-150.static.xvps.ne.jp/images/1a2b3c4d-5678-90ab-cdef-1234567890ab.jpg"
  ],
  "map": "https://www.google.com/maps/embed?pb=!1m34!1m12!1m3...",
  "checkpoints": [
	{ "name": "スタート地点", "lat": 33.5007, "lng": 129.8789 },
	{ "name": "絶景ポイント", "lat": 33.5370, "lng": 129.8950 },
	{ "name": "ゴール地点", "lat": 33.5555, "lng": 129.8463 }
  ]
}'
*/
