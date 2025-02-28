package domain

// ルート内の各チェックポイントを表現
type Checkpoint struct {
	Name string  `json:"name"` // チェックポイント名
	Lat  float64 `json:"lat"`  // 緯度
	Lng  float64 `json:"lng"`  // 経度
}
