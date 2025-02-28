package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// 全てのタグを取得して返す
func (h *Handler) GetTags(c *gin.Context) {
	tags, err := h.tagsService.GetTags()
	if err != nil {
		logger.LogError(err, "タグのサービスで死んだ")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tags":  tags,
		"error": "",
	})
}
