package service

import (
	// "database/sql"
	// "errors"
	// "time"

	// "github.com/google/uuid"
	// "github.com/sushigon-dev/sagaicle/internal/domain"
	repository "github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
	// "golang.org/x/crypto/bcrypt"
)

// TagsService はルート関連のユースケースを扱います。
type TagsService interface {
	GetTags() ([]string, error)
}

type routeService struct {
	repo repository.TagsRepository
}

// NewRouteService は新たな RouteService を生成します。
func NewRouteService(repo repository.TagsRepository) TagsService {
	return &routeService{repo: repo}
}

// GetTags はリポジトリから全てのタグを取得します。
func (s *routeService) GetTags() ([]string, error) {
	return s.repo.GetTags()
}
