package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
)

// ルートの作成リクエストを処理
func (h *Handler) CreateRoute(c *gin.Context) {
	var req struct {
		Title            string              `json:"title"`
		Description      string              `json:"description"`
		FullDescription  string              `json:"full_description"`
		Distance         float64             `json:"distance"`
		Time             int                 `json:"time"`
		Tags             []string            `json:"tags"`
		TotalCheckpoints int                 `json:"total_checkpoints"`
		Images           []string            `json:"images"`
		Map              string              `json:"map"`
		Checkpoints      []domain.Checkpoint `json:"checkpoints"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route := &domain.Route{
		Title:            req.Title,
		Description:      req.Description,
		FullDescription:  req.FullDescription,
		Distance:         req.Distance,
		Time:             req.Time,
		Tags:             req.Tags,
		TotalCheckpoints: req.TotalCheckpoints,
		Images:           req.Images,
		Map:              req.Map,
		Checkpoints:      req.Checkpoints,
	}

	// サービス層で入力検証、UUID発行、更新日時設定などが行われる
	if err := h.routesService.CreateRoute(route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"route_id":          route.ID.String(),
		"update_at":         route.UpdateAt.Format("2006/01/02"),
		"title":             route.Title,
		"description":       route.Description,
		"full_description":  route.FullDescription,
		"distance":          route.Distance,
		"time":              route.Time,
		"tags":              route.Tags,
		"total_checkpoints": route.TotalCheckpoints,
		"images":            route.Images,
		"map":               route.Map,
		"checkpoints":       route.Checkpoints,
		"error":             "",
	})
}

// パスパラメータからルートIDを取得し、ルート詳細を返す
func (h *Handler) GetRouteByID(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
		return
	}
	route, err := h.routesService.GetRouteByID(routeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"route_id":          route.ID.String(),
		"title":             route.Title,
		"description":       route.Description,
		"full_description":  route.FullDescription,
		"distance":          route.Distance,
		"time":              route.Time,
		"tags":              route.Tags,
		"total_checkpoints": route.TotalCheckpoints,
		"images":            route.Images,
		"map":               route.Map,
		"update_at":         route.UpdateAt.Format("2006/01/02"),
		"checkpoints":       route.Checkpoints,
		"error":             "",
	})
}

// 検索条件に基づいてルート一覧を取得
func (h *Handler) SearchRoutes(c *gin.Context) {
	/*
		var criteria domain.SearchCriteria
		if err := c.BindJSON(&criteria); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		routes, hitCount, err := h. routesService.SearchRoutes(&criteria)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"hit_count":     hitCount,
			"routes":        routes,
			"distance":      criteria.Distance,
			"time":          criteria.Time,
			"tags":          criteria.Tags,
			"search_option": criteria.SearchOption,
			"sort":          criteria.Sort,
			"limit":         criteria.Limit,
			"error":         "",
		})
	*/
}
