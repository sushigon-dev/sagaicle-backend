package repository

import (
// "github.com/sushigon-dev/sagaicle/internal/domain"
)

// タグに関する操作を抽象化したインターフェース
type TagsRepository interface {
	GetTags() ([]string, error)
}
