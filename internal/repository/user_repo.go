package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sushigon-dev/sagaicle/internal/domain"
)

// ユーザーに関する操作を抽象化したインターフェース
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByName(ctx context.Context, username string) (*domain.User, error)
}

type sqliteUserRepository struct {
	db *sqlx.DB
}

func NewSQLiteUserRepository(db *sqlx.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO Users (user_id, user_name, password_hash, mileage, total_distance)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.UserName, user.PasswordHash, user.Mileage, user.TotalDistance)
	return err
}

func (r *sqliteUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT user_id, user_name, password_hash, mileage, total_distance
		FROM Users WHERE user_id = ?
	`
	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqliteUserRepository) GetUserByName(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT user_id, user_name, password_hash, mileage, total_distance
		FROM Users WHERE user_name = ?
	`
	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		return nil, err
	}
	return &user, nil
}
