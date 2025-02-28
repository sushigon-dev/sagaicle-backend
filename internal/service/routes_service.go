package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

func NewRoutesService(repo sqlite.RoutesRepository) RoutesService {
	return &routesService{repo: repo}
}

// ルート作成のビジネスルールを実施し、リポジトリへ委譲
func (s *routesService) CreateRoute(route *domain.Route) error {
	// 入力検証（各フィールドの長さ、数値の範囲など）
	if len([]rune(route.Title)) == 0 || len([]rune(route.Title)) > 20 {
		err := errors.New("タイトルは1〜20文字である必要があります")
		logger.LogError(err, "タイトル文字長の検証エラー:"+route.Title)
		return err
	}
	if len([]rune(route.Description)) == 0 || len([]rune(route.Description)) > 60 {
		err := errors.New("説明は1〜60文字である必要があります")
		logger.LogError(err, "説明文字長の検証エラー"+route.Description)
		return err
	}
	if len([]rune(route.FullDescription)) == 0 || len([]rune(route.FullDescription)) > 200 {
		err := errors.New("詳細説明は1〜200文字である必要があります")
		logger.LogError(err, "詳細説明文字長の検証エラー"+route.FullDescription)
		return err
	}
	if route.Distance < 0 {
		err := errors.New("距離は0以上である必要があります")
		logger.LogError(err, "距離の検証エラー"+fmt.Sprint(route.Distance))
		return err
	}
	if route.Time < 0 {
		err := errors.New("所要時間は0以上である必要があります")
		logger.LogError(err, "所要時間の検証エラー"+fmt.Sprint(route.Time))
		return err
	}
	if len(route.Tags) > 20 {
		err := errors.New("タグは最大20個までです")
		logger.LogError(err, "タグ数の検証エラー"+fmt.Sprint(route.Tags))
		return err
	}
	for _, tag := range route.Tags {
		if len([]rune(tag)) < 1 || len([]rune(tag)) > 10 {
			err := errors.New("タグは1〜10文字である必要があります")
			logger.LogError(err, "各タグ文字長の検証エラー"+tag)
			return err
		}
	}
	if route.TotalCheckpoints < 1 || route.TotalCheckpoints > 20 {
		err := errors.New("チェックポイント数は1〜20である必要があります")
		logger.LogError(err, "チェックポイント数の検証エラー"+fmt.Sprint(route.TotalCheckpoints))
		return err
	}
	if len(route.Images) < 1 || len(route.Images) > 6 {
		err := errors.New("画像は1〜6件である必要があります")
		logger.LogError(err, "画像数の検証エラー"+fmt.Sprint(route.Images))
		return err
	}
	for _, image := range route.Images {
		// 文字数検証
		if len([]rune(image)) < 8 || len([]rune(image)) > 1023 {
			err := errors.New("画像URLは8〜1023文字である必要があります")
			logger.LogError(err, "画像URL文字長の検証エラー"+image)
			return err
		}

		// スキーム検証
		if image[:8] != "https://" {
			err := errors.New("画像は有効な形式である必要があります")
			logger.LogError(err, "画像URLのスキーム検証エラー"+image)
			return err
		}

		// ホスト検証
		if image[:18] != "https://github.com" &&
			image[:31] != "https://sushigon-dev.github.io/" {
			err := errors.New("画像は有効な形式である必要があります")
			logger.LogError(err, "画像URLのホスト検証エラー"+image)
			return err
		}

		// 拡張子検証
		if image[len(image)-4:] != ".png" && image[len(image)-4:] != ".jpg" &&
			image[len(image)-5:] != ".jpeg" && image[len(image)-5:] != ".webp" {
			err := errors.New("画像は有効な形式である必要があります")
			logger.LogError(err, "画像URLの拡張子検証エラー"+image)
			return err
		}
	}

	// ルート作成時に UUID と更新日時を設定
	route.ID = uuid.New()
	route.UpdateAt = time.Now()
	route.Likes = 0 // 初期状態はいいねなし

	return s.repo.CreateRoute(route)
}

// ルート ID からルート詳細情報を取得
func (s *routesService) GetRouteByID(id uuid.UUID) (*domain.Route, error) {
	return s.repo.GetRouteByID(id)
}
