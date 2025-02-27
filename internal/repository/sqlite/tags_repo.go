package repository

import (
	"github.com/jmoiron/sqlx"
)

type sqliteTagsRepository struct {
	db *sqlx.DB
}

func NewSQLiteTagsRepository(db *sqlx.DB) *sqliteTagsRepository {
	return &sqliteTagsRepository{db: db}
}

// 全てのタグを取得
func (r *sqliteTagsRepository) GetTags() ([]string, error) {
	var tags []string
	query := `SELECT tag FROM tags;`
	if err := r.db.Select(&tags, query); err != nil {
		return nil, err
	}
	return tags, nil
}
