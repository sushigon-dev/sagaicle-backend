package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ユーザーが訪問したチェックポイント一覧を取得
// ここでは、サービス層でそのロジックを実装しているものと仮定
func (h *Handler) GetVisitedCheckpoints(c *gin.Context) {
	/*
		routeIDStr := c.Param("route_id")
		routeID, err := uuid.Parse(routeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
			return
		}
		userIDStr := c.GetHeader("X-User-ID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

			// サービス層に実装があると仮定
			visited, err := h.checkpointsService.VisitCheckpoint(userID, routeID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
		return
	}
	cpIndexStr := c.Param("checkpoint_index")
	cpIndex, err := strconv.Atoi(cpIndexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid checkpoint index"})
		return
	}
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	if err := h.checkpointsService.VisitCheckpoint(
		userID, routeID, cpIndex,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// レスポンスとしては、更新後の状態
	c.JSON(http.StatusOK, gin.H{
		"route_id":         routeID.String(),
		"checkpoint_index": cpIndex,
		"error":            "",
	})
}
