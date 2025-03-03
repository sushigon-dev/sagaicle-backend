package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/utils/errors"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// ユーザーが訪問したチェックポイント一覧を取得
// ここでは、サービス層でそのロジックを実装しているものと仮定
func (h *Handler) GetVisitedCheckpoints(c *gin.Context) {
	/*
				routeIDStr := c.Param("route_id")
				routeID, err := uuid.Parse(routeIDStr)
				if err != nil {
		            logger.LogError(err, "route_idのパースに失敗")
					c.JSON(http.StatusBadRequest, gin.H{"error": })
					return
				}
				userIDStr := c.GetHeader("X-User-ID")
				userID, err := uuid.Parse(userIDStr)
				if err != nil {
		            logger.LogError(err, "X-User-IDのパースに失敗")
					c.JSON(http.StatusUnauthorized, gin.H{"error": })
					return
				}

					// サービス層に実装があると仮定
					visited, err := h.checkpointsService.VisitCheckpoint(userID, routeID)
					if err != nil {
		                logger.LogError(err, "チェックポイント訪問記録の取得に失敗")
						c.JSON(http.StatusInternalServerError, gin.H{"error": })
						return
					}
					c.JSON(http.StatusOK, gin.H{
						"route_id":            routeID.String(),
						"visited_checkpoints": visited,
						"error":               "",
					})
	*/
}

// 指定のチェックポイント訪問記録を追加
func (h *Handler) VisitCheckpoint(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		logger.LogError(err, "route_idのパースに失敗")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.InvalidFormat})
		return
	}

	cpIndexStr := c.Param("checkpoint_index")
	cpIndex, err := strconv.Atoi(cpIndexStr)
	if err != nil {
		logger.LogError(err, "checkpoint_indexのパースに失敗")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.InvalidFormat})
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.LogError(err, "X-User-IDのパースに失敗")
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.InvalidFormat})
		return
	}

	if err := h.checkpointsService.VisitCheckpoint(
		userID, routeID, cpIndex,
	); err != nil {
		logger.LogError(err, "チェックポイント訪問記録の追加に失敗")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.InternalServer})
		return
	}

	// レスポンスとしては、更新後の状態
	c.JSON(http.StatusOK, gin.H{
		"route_id":         routeID.String(),
		"checkpoint_index": cpIndex,
		"error":            "",
	})
}
