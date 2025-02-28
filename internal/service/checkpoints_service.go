package service

import (
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
)

func NewCheckpointsService(repo sqlite.CheckpointsRepository) CheckpointsService {
	return &checkpointsService{repo: repo}
}

// ユーザーが特定のチェックポイントを訪問した記録を追加
func (s *checkpointsService) VisitCheckpoint(userID, routeID uuid.UUID, checkpointIndex int) error {
	return s.repo.VisitCheckpoint(userID, routeID, checkpointIndex)
}
