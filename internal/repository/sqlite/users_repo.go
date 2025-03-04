package sqlite

import (
	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// 新規ユーザーを登録
func (r *SQLiteRepository) CreateUser(user *domain.User) error {
	query := `
        INSERT INTO users (
            user_id, user_name, hashed_password, mileage, total_distance
        ) VALUES (
            ?, ?, ?, ?, ?
        );
    `
	if _, err := r.db.Exec(query, user.ID.String(),
		user.UserName, user.PasswordHash, user.Mileage, user.TotalDistance,
	); err != nil {
		logger.LogError(err, "ユーザー登録に失敗")
		return err
	}

	return nil
}

// ユーザー ID によるユーザー情報取得
func (r *SQLiteRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	query := `
        SELECT user_id, user_name, hashed_password, mileage, total_distance
        FROM users
        WHERE user_id = ?;
    `
	var row struct {
		UserID         string  `db:"user_id"`
		UserName       string  `db:"user_name"`
		HashedPassword string  `db:"hashed_password"`
		Mileage        float64 `db:"mileage"`
		TotalDistance  float64 `db:"total_distance"`
	}
	if err := r.db.Get(&row, query, id.String()); err != nil {
		logger.LogError(err, "ID によるユーザー情報の取得に失敗")
		return nil, err
	}

	user := &domain.User{
		ID:            id,
		UserName:      row.UserName,
		PasswordHash:  row.HashedPassword,
		Mileage:       row.Mileage,
		TotalDistance: row.TotalDistance,
	}
	return user, nil
}

// ユーザー名からユーザー情報を取得
func (r *SQLiteRepository) GetUserByUsername(username string) (*domain.User, error) {
	query := `
        SELECT user_id, user_name, hashed_password, mileage, total_distance
        FROM users
        WHERE user_name = ?;
    `
	var row struct {
		UserID         string  `db:"user_id"`
		UserName       string  `db:"user_name"`
		HashedPassword string  `db:"hashed_password"`
		Mileage        float64 `db:"mileage"`
		TotalDistance  float64 `db:"total_distance"`
	}
	if err := r.db.Get(&row, query, username); err != nil {
		logger.LogError(err, "ユーザー名からのユーザー情報取得に失敗")
		return nil, err
	}

	uid, err := uuid.Parse(row.UserID)
	if err != nil {
		logger.LogError(err, "User ID のパースに失敗")
		return nil, err
	}

	user := &domain.User{
		ID:            uid,
		UserName:      row.UserName,
		PasswordHash:  row.HashedPassword,
		Mileage:       row.Mileage,
		TotalDistance: row.TotalDistance,
	}
	return user, nil
}

// ユーザーが取得したバッジ付きルートの記録を追加
func (r *SQLiteRepository) AddBadgedRoute(userID uuid.UUID, route domain.BadgedRoute) error {
	query := `
        INSERT INTO badges (user_id, route_id, title)
        VALUES (?, ?, ?);
    `
	if _, err := r.db.Exec(query, userID.String(), route.ID, route.Title); err != nil {
		logger.LogError(err, "バッジ付きルートの記録追加に失敗")
		return err
	}

	return nil
}
