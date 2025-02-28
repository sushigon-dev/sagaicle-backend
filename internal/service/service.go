package service

import (
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
)

// チェックポイント関連のユースケースを扱う
type CheckpointsService interface {
	VisitCheckpoint(userID, routeID uuid.UUID, checkpointIndex int) error
}

type checkpointsService struct {
	repo sqlite.CheckpointsRepository
}

// いいね関連のユースケースを扱う
type LikesService interface {
	LikeRoute(userID, routeID uuid.UUID) error
	DislikeRoute(userID, routeID uuid.UUID) error
	IsLiked(userID, routeID uuid.UUID) (bool, int, error)
}

type likesService struct {
	repo sqlite.LikesRepository
}

// ルート関連のユースケースを扱う
type RoutesService interface {
	CreateRoute(route *domain.Route) error
	GetRouteByID(id uuid.UUID) (*domain.Route, error)
	SearchRoutes(criteria *domain.SearchCriteria) ([]domain.RouteSummary, int, error)
}

type routesService struct {
	repo sqlite.RoutesRepository
}

// タグ関連のユースケースを扱う
type TagsService interface {
	GetTags() ([]string, error)
}

type tagsService struct {
	repo sqlite.TagsRepository
}

// ユーザー関連のユースケースを扱う
type UsersService interface {
	Register(userName, password string) (*domain.User, error)
	Login(userName, password string) (*domain.User, error)
	GetUserProfile(userID uuid.UUID) (*domain.User, error)
}

type usersService struct {
	repo sqlite.UsersRepository
}
