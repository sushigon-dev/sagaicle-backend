package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sushigon-dev/sagaicle/utils/errors"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// 全てのタグを取得して返す
func (h *Handler) GetTags(c *gin.Context) {
	tags, err := h.tagsService.GetTags()
	if err != nil {
		logger.Error(err, "タグの取得に失敗")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.InternalServer})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tags":  tags,
		"error": "",
	})
}
