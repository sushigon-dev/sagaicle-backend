package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sushigon-dev/sagaicle/internal/handler"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
	"github.com/sushigon-dev/sagaicle/internal/service"
)

func main() {
	// db の設定
	db := configDB()
	defer db.Close()

	// リポジトリ、サービス、ハンドラーの初期化
	repo := sqlite.NewSQLiteRepository(db)
	service := service.NewRouteService(repo)
	h := handler.NewHandler(service)

	// Gin ルーターの作成
	router := gin.Default()

	// CORS ミドルウェアの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "https://sushigon-dev.github.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ルートの登録
	h.RegisterRoutes(router)

	// サーバー起動
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
	}
}
