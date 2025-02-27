package handler

import (
	"net/http"
	// "strconv"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	// "github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/internal/service"
)

// Handler はサービス層を利用するハンドラー層のエントリーポイントです。
type Handler struct {
	routeService service.TagsService
}

// NewHandler は新たなハンドラーを生成します。
func NewHandler(routeService service.TagsService) *Handler {
	return &Handler{
		routeService: routeService,
	}
}

// RegisterRoutes は各エンドポイントをルーターに登録します。
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// タグ
	router.GET("/api/tags", h.GetTags)
}

// GetTags は全てのタグを取得して返します。
func (h *Handler) GetTags(c *gin.Context) {
	tags, err := h.routeService.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tags":  tags,
		"error": "",
	})
}
