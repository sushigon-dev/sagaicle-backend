package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite ドライバ
	"github.com/sushigon-dev/sagaicle/internal/handler"
	repository "github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
	"github.com/sushigon-dev/sagaicle/internal/service"
)

func main() {
	// sqlx を利用して SQLite データベースに接続
	db, err := sqlx.Open("sqlite3", "./sagaicle.db")
	if err != nil {
		log.Fatalf("SQLite データベースのオープンに失敗しました: %v", err)
	}
	defer db.Close()

	// データベース接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("SQLite データベースへの接続確認に失敗しました: %v", err)
	}

	// 必要なテーブルがなければ作成する
	if err := createTables(db); err != nil {
		log.Fatalf("テーブル作成に失敗しました: %v", err)
	}

	// リポジトリ、サービス、ハンドラーの初期化
	tagsRepo := repository.NewSQLiteTagsRepository(db)
	tagsService := service.NewRouteService(tagsRepo)
	h := handler.NewHandler(tagsService)

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

// createTables は必要なテーブルを作成します。
// ここでは tags テーブルを例として作成しています。
func createTables(db *sqlx.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS tags (
        tag TEXT PRIMARY KEY,
        CHECK (LENGTH(tag) BETWEEN 1 AND 10)
    );
	`
	_, err := db.Exec(schema)
	return err
}
