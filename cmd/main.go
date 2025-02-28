package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sushigon-dev/sagaicle/internal/handler"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
	"github.com/sushigon-dev/sagaicle/internal/service"
	"github.com/sushigon-dev/sagaicle/utils/config"
)

func main() {
	// ログ出力先を標準出力に設定
	log.SetOutput(os.Stdout)

	// db の設定
	db := config.DB()
	defer db.Close()

	// リポジトリの初期化
	repo := sqlite.NewSQLiteRepository(db)

	// サービスの初期化
	cs := service.NewCheckpointsService(repo)
	ls := service.NewLikesService(repo)
	rs := service.NewRoutesService(repo)
	ts := service.NewTagsService(repo)
	us := service.NewUserService(repo)

	// ハンドラーの初期化
	h := handler.NewHandler(cs, ls, rs, ts, us)

	// Gin ルーターの作成
	r := gin.Default()

	// 標準のロガーとリカバリミドルウェアを設定
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS ミドルウェアの設定
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// origin が 次のいずれかで始まれば許可する
			return strings.HasPrefix(origin, "http://localhost") ||
				strings.HasPrefix(origin, "https://sushigon-dev.github.io") ||
				strings.HasPrefix(origin, "https://sagaicle-frontend-btwp.vercel.app")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ルートの登録
	h.RegisterRoutes(r)

	// サーバー起動
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
	}
}
