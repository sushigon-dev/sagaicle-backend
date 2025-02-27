package service

import (
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
)

func NewLikesService(repo sqlite.LikesRepository) LikesService {
	return &likesService{repo: repo}
}

// 指定ユーザーによる「いいね」操作を実行
func (s *likesService) LikeRoute(userID, routeID uuid.UUID) error {
	return s.repo.LikeRoute(userID, routeID)
}

// 指定ユーザーによる「いいね」削除操作を実行
func (s *likesService) DislikeRoute(userID, routeID uuid.UUID) error {
	return s.repo.DislikeRoute(userID, routeID)
}

// ユーザーがルートをいいねしているかどうかと総いいね数を返す
func (s *likesService) IsLiked(userID, routeID uuid.UUID) (bool, int, error) {
	return s.repo.IsLiked(userID, routeID)
}
