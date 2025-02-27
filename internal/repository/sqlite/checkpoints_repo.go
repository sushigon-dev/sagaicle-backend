package sqlite

import (
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
)

// CreateCheckpoints は指定したルート ID に対するチェックポイント群を登録します。
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
		return err
	}
	for idx, cp := range checkpoints {
		if _, err := tx.Exec(query, routeID.String(), idx, cp.Name, cp.Lat, cp.Lng); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// GetCheckpointsByRouteID はルート ID に紐づくチェックポイント群を取得します。
func (r *SQLiteRepository) GetCheckpointsByRouteID(routeID uuid.UUID) ([]domain.Checkpoint, error) {
	query := `
        SELECT name, lat, lng
        FROM checkpoints
        WHERE route_id = ?
        ORDER BY checkpoint_index ASC;
    `
	var checkpoints []domain.Checkpoint
	if err := r.db.Select(&checkpoints, query, routeID.String()); err != nil {
		return nil, err
	}
	return checkpoints, nil
}
