package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 新規ユーザー登録のリクエストを処理
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.usersService.Register(req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 登録成功後、Cookie の設定等も行えます（ここではシンプルにユーザー名のみ返す）
	c.JSON(http.StatusCreated, gin.H{
		"user_name": user.UserName,
		"error":     "",
	})
}

// ログイン認証を行い、成功時はセッション情報（例：トークン）を返す
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.usersService.Login(req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// ここでトークン発行や Cookie 設定を行うのが一般的です
	c.JSON(http.StatusOK, gin.H{
		"user_name": user.UserName,
		"error":     "",
	})
}

// 認証済みユーザーの情報を返す
func (h *Handler) Whoami(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, err := h.usersService.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":   user.ID.String(),
		"user_name": user.UserName,
		"error":     "",
	})
}

// ログアウト処理（例：トークンの無効化）を行う
func (h *Handler) Logout(c *gin.Context) {
	/*
		// 実際にはトークン無効化や Cookie のクリア等の処理を行う
		c.JSON(http.StatusOK, gin.H{
			"token": "",
			"error": "",
		})
	*/
}
