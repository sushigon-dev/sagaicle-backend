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
type RouteRepository interface {
	CreateRoute(route *domain.Route) error
	GetRouteByID(id uuid.UUID) (*domain.Route, error)
}

// sqlx.DB を利用したリポジトリ実装
type SQLiteRepository struct {
	db *sqlx.DB
}

// 新たなリポジトリインスタンスを生成
func NewSQLiteRepository(db *sqlx.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}
