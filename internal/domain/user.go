package model

import (
	"github.com/google/uuid"
)

// ユーザー情報を表現
type User struct {
	ID            uuid.UUID `json:"id"`             // ユーザーのユニークID
	UserName      string    `json:"user_name"`      // ユーザー名
	PasswordHash  string    `json:"password_hash"`  // パスワードのハッシュ値（実際のパスワードは保存しない）
	Mileage       float64   `json:"mileage"`        // 累積マイレージ
	TotalDistance float64   `json:"total_distance"` // 走破した総距離

	BadgedRoutes []BadgedRoute `json:"badged_routes"` // バッジが付与されたルート一覧
	LikedRoutes  []LikedRoute  `json:"liked_routes"`  // いいね済みルート一覧
}

// ユーザーのバッジ付きルートの簡易情報
type BadgedRoute struct {
	ID    string `json:"id"`       // ルートID（文字列として管理する場合もある）
	Title string `json:"title"   ` // タイトル
}

// ユーザーのいいね済みルートの簡易情報
type LikedRoute struct {
	ID    uuid.UUID `json:"id"`       // ルートのユニークID
	Title string    `json:"title"   ` // タイトル
}
