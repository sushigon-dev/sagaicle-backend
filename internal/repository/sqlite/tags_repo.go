package sqlite

// 全てのタグを取得
func (r *SQLiteRepository) GetTags() ([]string, error) {
	var tags []string
	query := `SELECT tag FROM tags;`
	if err := r.db.Select(&tags, query); err != nil {
		return nil, err
	}
	return tags, nil
}
