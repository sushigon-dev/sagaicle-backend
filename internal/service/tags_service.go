package service

import (
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
)

func NewTagsService(repo sqlite.TagsRepository) TagsService {
	return &tagsService{repo: repo}
}

// リポジトリから全てのタグを取得
func (s *tagsService) GetTags() ([]string, error) {
	return s.repo.GetTags()
}
