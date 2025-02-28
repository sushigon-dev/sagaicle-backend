package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/utils/errors"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// 指定ユーザーがルートをいいねしているかを確認
func (h *Handler) IsLiked(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		logger.LogError(err, "route_idのパースに失敗")
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

	isLiked, likes, err := h.likesService.IsLiked(userID, routeID)
	if err != nil {
		logger.LogError(err, "いいねの確認に失敗")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.InternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"route_id": routeID.String(),
		"is_liked": isLiked,
		"likes":    likes,
		"error":    "",
	})
}

// ルートへの「いいね」を追加
func (h *Handler) LikeRoute(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		logger.LogError(err, "route_idのパースに失敗")
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

	if err := h.likesService.LikeRoute(userID, routeID); err != nil {
		logger.LogError(err, "いいねの追加に失敗")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.InternalServer})
		return
	}

	_, likes, err := h.likesService.IsLiked(userID, routeID)
	if err != nil {
		logger.LogError(err, "いいねの確認に失敗")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.InternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"route_id": routeID.String(),
		"likes":    likes,
		"error":    "",
	})
}

// ルートから「いいね」を削除
func (h *Handler) DislikeRoute(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		logger.LogError(err, "route_idのパースに失敗")
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

	if err := h.likesService.DislikeRoute(userID, routeID); err != nil {
		logger.LogError(err, "いいねの削除に失敗")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.InternalServer})
		return
	}

	_, likes, err := h.likesService.IsLiked(userID, routeID)
	if err != nil {
		logger.LogError(err, "いいねの確認に失敗")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.InternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"route_id": routeID.String(),
		"likes":    likes,
		"error":    "",
	})
}
