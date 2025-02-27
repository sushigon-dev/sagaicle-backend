package sqlite

import (
	// "database/sql"
	// "encoding/json"
	// "time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sushigon-dev/sagaicle/internal/domain"
)

// タグに関する操作を抽象化したインターフェース
type TagsRepository interface {
	GetTags() ([]string, error)
}

// ルートに関する操作を抽象化したインターフェース
type RoutesRepository interface {
	CreateRoute(route *domain.Route) error
	GetRouteByID(id uuid.UUID) (*domain.Route, error)
}

// チェックポイントに関する操作を抽象化したインターフェース
type CheckpointsRepository interface {
	CreateCheckpoints(routeID uuid.UUID, checkpoints []domain.Checkpoint) error
	GetCheckpointsByRouteID(routeID uuid.UUID) ([]domain.Checkpoint, error)
}

// ユーザーに関する操作を抽象化したインターフェース
type UsersRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id uuid.UUID) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}

// sqlx.DB を利用したリポジトリ実装
type SQLiteRepository struct {
	db *sqlx.DB
}

// 新たなリポジトリインスタンスを生成
func NewSQLiteRepository(db *sqlx.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}
