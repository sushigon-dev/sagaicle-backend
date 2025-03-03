package sqlite

import (
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// 指定したルート ID に対するチェックポイント群を登録
func (r *SQLiteRepository) CreateCheckpoints(routeID uuid.UUID, checkpoints []domain.Checkpoint) error {
	// checkpoints テーブルは (route_id, checkpoint_index, name, lat, lng) を持つと想定
	query := `
        INSERT INTO checkpoints (
            route_id, checkpoint_index, name, lat, lng
        ) VALUES (
            ?, ?, ?, ?, ?
        );
    `
	tx, err := r.db.Beginx()
	if err != nil {
		logger.LogError(err, "トランザクションの開始に失敗")
		return err
	}

	for idx, cp := range checkpoints {
		if _, err := tx.Exec(query, routeID.String(), idx, cp.Name, cp.Lat, cp.Lng); err != nil {
			tx.Rollback()
			logger.LogError(err, "チェックポイントの登録に失敗")
			return err
		}
	}

	return tx.Commit()
}

// ルート ID に紐づくチェックポイント群を取得
func (r *SQLiteRepository) GetCheckpointsByRouteID(routeID uuid.UUID) ([]domain.Checkpoint, error) {
	query := `
        SELECT name, lat, lng
        FROM checkpoints
        WHERE route_id = ?
        ORDER BY checkpoint_index ASC;
    `
	var checkpoints []domain.Checkpoint
	if err := r.db.Select(&checkpoints, query, routeID.String()); err != nil {
		logger.LogError(err, "チェックポイントの取得に失敗")
		return nil, err
	}

	return checkpoints, nil
}

// ユーザーが指定ルートのチェックポイント（インデックス指定）を訪問した記録を追加
func (r *SQLiteRepository) VisitCheckpoint(userID, routeID uuid.UUID, checkpointIndex int) error {
	query := `
        INSERT INTO visited_checkpoints (user_id, route_id, checkpoint_index)
        VALUES (?, ?, ?);
    `
	_, err := r.db.Exec(query, userID.String(), routeID.String(), checkpointIndex)
	if err != nil {
		logger.LogError(err, "チェックポイントの訪問記録の追加に失敗")
	}

	return err
}
