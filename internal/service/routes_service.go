package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
)

func NewRoutesService(repo sqlite.RoutesRepository) RoutesService {
	return &routesService{repo: repo}
}

// ルート作成のビジネスルールを実施し、リポジトリへ委譲
func (s *routesService) CreateRoute(route *domain.Route) error {
	// 入力検証（各フィールドの長さ、数値の範囲など）
	if len(route.Title) == 0 || len(route.Title) > 20 {
		return errors.New("タイトルは1〜20文字である必要があります")
	}
	if len(route.Description) == 0 || len(route.Description) > 60 {
		return errors.New("説明は1〜60文字である必要があります")
	}
	if len(route.FullDescription) == 0 || len(route.FullDescription) > 200 {
		return errors.New("詳細説明は1〜200文字である必要があります")
	}
	if route.Distance < 0 {
		return errors.New("距離は0以上である必要があります")
	}
	if route.Time < 0 {
		return errors.New("所要時間は0以上である必要があります")
	}
	if len(route.Tags) > 20 {
		return errors.New("タグは最大20個までです")
	}
	if route.TotalCheckpoints < 1 || route.TotalCheckpoints > 20 {
		return errors.New("チェックポイント数は1〜20である必要があります")
	}
	if len(route.Images) < 1 || len(route.Images) > 6 {
		return errors.New("画像は1〜6件である必要があります")
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
