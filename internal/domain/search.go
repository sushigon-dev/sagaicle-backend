package domain

// 距離などの浮動小数点数の範囲を表現
type RangeFloat struct {
	Min float64 `json:"min"` // 最小値
	Max float64 `json:"max"` // 最大値
}

// 時間などの整数の範囲を表現
type RangeInt struct {
	Min int `json:"min"` // 最小値
	Max int `json:"max"` // 最大値
}

// ソート条件を表現
type Sort struct {
	Key   string `json:"key"`   // ソートキー "distance", "time", "likes" など
	Order string `json:"order"` // "asc" または "desc"
}

// ルート検索の際のパラメータをまとめた値オブジェクト
type SearchCriteria struct {
	Distance     RangeFloat
	Time         RangeInt
	Tags         []string // タグ一覧
	SearchOption []string // "AND", "OR", "NOT" などの検索オプション
	Sort         Sort
	Limit        int // 最大件数
}
