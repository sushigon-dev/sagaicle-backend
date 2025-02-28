package sqlite

import (
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// ユーザーがルートに「いいね」した記録を追加、ルートの likes カウントを更新
func (r *SQLiteRepository) LikeRoute(userID, routeID uuid.UUID) error {
	// likes テーブルに挿入
	query := `
        INSERT INTO likes (user_id, route_id)
        VALUES (?, ?);
    `
	if _, err := r.db.Exec(query, userID.String(), routeID.String()); err != nil {
		logger.LogError(err, "いいねの追加に失敗")
		return err
	}

	// routes テーブルの likes カウントをインクリメント
	updateQuery := `
        UPDATE routes SET likes = likes + 1 WHERE id = ?;
    `
	if _, err := r.db.Exec(updateQuery, routeID.String()); err != nil {
		logger.LogError(err, "いいね数の更新に失敗")
		return err
	}

	return nil
}

// ユーザーの「いいね」を削除し、ルートの likes カウントを更新
func (r *SQLiteRepository) DislikeRoute(userID, routeID uuid.UUID) error {
	query := `
        DELETE FROM likes
        WHERE user_id = ? AND route_id = ?;
    `
	res, err := r.db.Exec(query, userID.String(), routeID.String())
	if err != nil {
		logger.LogError(err, "いいねの削除に失敗")
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		logger.LogError(err, "いいねの削除に失敗")
		return err
	}

	if affected > 0 {
		updateQuery := `
            UPDATE routes SET likes = likes - 1 WHERE id = ?;
        `
		_, err = r.db.Exec(updateQuery, routeID.String())
		if err != nil {
			logger.LogError(err, "いいね数の更新に失敗")
		}
		return err
	}

	return nil
}

// 指定ユーザーがルートをいいねしているかどうかと、総いいね数を返す
func (r *SQLiteRepository) IsLiked(userID, routeID uuid.UUID) (bool, int, error) {
	query := `
        SELECT COUNT(1) FROM likes
        WHERE user_id = ? AND route_id = ?;
    `
	var count int
	if err := r.db.Get(&count, query, userID.String(), routeID.String()); err != nil {
		logger.LogError(err, "いいねの取得に失敗")
		return false, 0, err
	}

	// ルートの総いいね数を取得
	queryTotal := `
        SELECT likes FROM routes WHERE id = ?;
    `
	var likes int
	if err := r.db.Get(&likes, queryTotal, routeID.String()); err != nil {
		logger.LogError(err, "いいね数の取得に失敗")
		return false, 0, err
	}

	return count > 0, likes, nil
}
