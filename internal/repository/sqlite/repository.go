package sqlite

import (
	// "github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/jmoiron/sqlx"
)

// タグに関する操作を抽象化したインターフェース
type TagsRepository interface {
	GetTags() ([]string, error)
}

// SQLiteRepository は sqlx.DB を利用したリポジトリ実装です。
type SQLiteRepository struct {
	db *sqlx.DB
}

// NewSQLiteRepository は新たなリポジトリインスタンスを生成します。
func NewSQLiteRepository(db *sqlx.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}
