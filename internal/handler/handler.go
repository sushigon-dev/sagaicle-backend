package handler

import (
	// "strconv"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	// "github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/internal/service"
)

// サービス層を利用するハンドラー層のエントリーポイント
type Handler struct {
	tagsService service.TagsService
}

// 新たなハンドラーを生成
func NewHandler(tagsService service.TagsService) *Handler {
	return &Handler{
		tagsService: tagsService,
	}
}

// 各エンドポイントをルーターに登録
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// タグ
	router.GET("/api/tags", h.GetTags)
}
