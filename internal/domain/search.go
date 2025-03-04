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
	Distance     RangeFloat `json:"distance"`      // 距離の範囲
	Time         RangeInt   `json:"time"`          // 時間の範囲
	Tags         []string   `json:"tags"`          // 検索に使うタグ名
	SearchOption string     `json:"search_option"` // タグの検索方法 "AND", "OR", "NOT"
	Sort         Sort       `json:"sort"`          // 並び替え（キーと順序の）指定
	Limit        int        `json:"limit"`         // 取得する最大件数
}
