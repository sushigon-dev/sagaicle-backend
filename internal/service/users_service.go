package service

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/internal/repository/sqlite"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(repo sqlite.UsersRepository) UsersService {
	return &usersService{repo: repo}
}

// ユーザー登録のビジネスルールを実装
// 既に存在するユーザー名のチェックや、パスワードのハッシュ化を行う
func (s *usersService) Register(userName, password string) (*domain.User, error) {
	if len(userName) == 0 || len(userName) > 32 {
		return nil, errors.New("ユーザー名は1〜32文字である必要があります")
	}
	if len(password) == 0 || len(password) > 32 {
		return nil, errors.New("パスワードは1〜32文字である必要があります")
	}

	// 既存ユーザーの存在確認（存在する場合はエラーを返す）
	if _, err := s.repo.GetUserByUsername(userName); err == nil {
		return nil, errors.New("既に存在するユーザー名です")
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:            uuid.New(),
		UserName:      userName,
		PasswordHash:  string(hashedPassword),
		Mileage:       0,
		TotalDistance: 0,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// ユーザー認証を行う
// パスワードの比較により認証の正否を判定
func (s *usersService) Login(userName, password string) (*domain.User, error) {
	user, err := s.repo.GetUserByUsername(userName)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("認証に失敗しました")
	}
	return user, nil
}

// ユーザー ID からユーザープロフィール情報を取得
func (s *usersService) GetUserProfile(userID uuid.UUID) (*domain.User, error) {
	return s.repo.GetUserByID(userID)
}
