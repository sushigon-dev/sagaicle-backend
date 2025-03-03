package sqlite

import (
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// 全てのタグを取得
func (r *SQLiteRepository) GetTags() ([]string, error) {
	var tags []string
	query := `SELECT tag FROM tags;`
	if err := r.db.Select(&tags, query); err != nil {
		logger.LogError(err, "タグの取得に失敗")
		return nil, err
	}
	return tags, nil
}
