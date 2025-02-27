package model

import (
	"time"

	"github.com/google/uuid"
)

// ルート詳細情報を表現
type Route struct {
	ID               uuid.UUID    `json:"id"`                // ルートのユニークID
	Title            string       `json:"title"`             //  タイトル (1〜20文字)
	Description      string       `json:"description"`       // 短い説明 (1〜60文字)
	FullDescription  string       `json:"full_description"`  // 詳細説明 (1〜200文字)
	Distance         float64      `json:"distance"`          // 距離（0以上の値）
	Time             int          `json:"time"`              // 時間（0以上の値、単位は分など）
	Tags             []string     `json:"tags"`              // ルートに関連付けられたタグ
	TotalCheckpoints int          `json:"total_checkpoints"` // チェックポイントの総数 (1〜20)
	Images           []string     `json:"images"`            // ルートの画像URL (1〜6件)
	Map              string       `json:"map"`               // マップのURL
	Checkpoints      []Checkpoint `json:"checkpoints"`       // ルートに含まれるチェックポイントの一覧
	Likes            int          `json:"likes"`             // いいねの総数
	UpdateAt         time.Time    `json:"update_at"`         // 更新日時（YYYY/MM/DD）
}
