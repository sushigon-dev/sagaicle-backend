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
	checkpointsService service.CheckpointsService
	likesService       service.LikesService
	routesService      service.RoutesService
	tagsService        service.TagsService
	usersService       service.UsersService
}

// 新たなハンドラーを生成
func NewHandler(checkpointsService service.CheckpointsService, likesService service.LikesService, routesService service.RoutesService, tagsService service.TagsService, usersService service.UsersService) *Handler {
	return &Handler{
		checkpointsService: checkpointsService,
		likesService:       likesService,
		routesService:      routesService,
		tagsService:        tagsService,
		usersService:       usersService,
	}
}

// 各エンドポイントをルーターに登録
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// タグ
	router.GET("/api/tags", h.GetTags)

	// ルート
	router.POST("/api/search", h.SearchRoutes)
	router.POST("/api/routes", h.CreateRoute)
	router.GET("/api/routes/:route_id", h.GetRouteByID)

	// いいね
	router.GET("/api/routes/:route_id/like", h.IsLiked)
	router.POST("/api/routes/:route_id/like", h.LikeRoute)
	router.DELETE("/api/routes/:route_id/like", h.DislikeRoute)

	// チェックポイント
	router.GET("/api/routes/:route_id/checkpoints", h.GetVisitedCheckpoints)
	router.POST("/api/routes/:route_id/checkpoints/:checkpoint_index/visit", h.VisitCheckpoint)

	// ユーザー（認証系）
	router.POST("/api/auth/register", h.Register)
	router.POST("/api/auth/login", h.Login)
	router.GET("/api/auth/me", h.Whoami)
	router.DELETE("/api/auth/logout", h.Logout)
}
