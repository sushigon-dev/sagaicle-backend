package sqlite

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sushigon-dev/sagaicle/internal/domain"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

// ルートの作成と、それに紐づくチェックポイントの登録
func (r *SQLiteRepository) CreateRoute(route *domain.Route) error {
	// update_at は "YYYY/MM/DD" 形式に変換
	updateAt := route.UpdateAt.Format("2006/01/02")

	// ルートテーブルへの INSERT クエリ（テーブルが違うものは後述）
	query := `
        INSERT INTO routes (
            id, title, description, full_description, distance, time, likes, image, update_at, total_checkpoints, map
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        );
    `
	// サマリー用に画像は配列の最初の要素を利用
	var image string
	if len(route.Images) > 0 {
		image = route.Images[0]
	}

	var err error
	_, err = r.db.Exec(query,
		route.ID.String(), route.Title, route.Description, route.FullDescription,
		route.Distance, route.Time, route.Likes, image,
		updateAt, route.TotalCheckpoints, route.Map,
	)
	if err != nil {
		logger.LogError(err, "ルートの登録に失敗"+fmt.Sprint(query,
			route.ID.String(), route.Title, route.Description, route.FullDescription,
			route.Distance, route.Time, route.Likes, image, updateAt,
			route.TotalCheckpoints, route.Map),
		)
		return err
	}

	// タグの登録：Route_Tagsテーブルに対して、ルートIDとタグ名を登録する
	for _, tag := range route.Tags {
		tagQuery := `INSERT INTO route_tags (route_id, tag_name) VALUES (?, ?);`
		_, err = r.db.Exec(tagQuery, route.ID.String(), tag)
		if err != nil {
			logger.LogError(err, "タグの登録に失敗: "+
				fmt.Sprint(tagQuery, route.ID.String(), tag))
			return err
		}
	}

	// 画像の登録：Route_Imagesテーブルに対して、ルートIDと画像URIを登録する
	for _, imageURI := range route.Images {
		imageQuery := `INSERT INTO route_images (route_id, image_uri) VALUES (?, ?);`
		_, err = r.db.Exec(imageQuery, route.ID.String(), imageURI)
		if err != nil {
			logger.LogError(err, "画像の登録に失敗: "+imageURI)
			return err
		}
	}

	// チェックポイントの登録
	return r.CreateCheckpoints(route.ID, route.Checkpoints)
}

// 指定したルート ID の詳細情報を取得
func (r *SQLiteRepository) GetRouteByID(id uuid.UUID) (*domain.Route, error) {
	query := `
        SELECT id, title, description, full_description, distance, time, tags, likes, image, update_at, total_checkpoints, images, map
        FROM routes
        WHERE id = ?;
    `
	var row struct {
		ID               string  `db:"id"`
		Title            string  `db:"title"`
		Description      string  `db:"description"`
		FullDescription  string  `db:"full_description"`
		Distance         float64 `db:"distance"`
		Time             int     `db:"time"`
		Tags             string  `db:"tags"`
		Likes            int     `db:"likes"`
		Image            string  `db:"image"`
		UpdateAt         string  `db:"update_at"`
		TotalCheckpoints int     `db:"total_checkpoints"`
		Images           string  `db:"images"`
		Map              string  `db:"map"`
	}

	if err := r.db.Get(&row, query, id.String()); err != nil {
		logger.LogError(err, "ルートの取得に失敗"+fmt.Sprint(query, id.String()))
		return nil, err
	}

	// JSON フィールドのデコード
	var tags []string
	if err := json.Unmarshal([]byte(row.Tags), &tags); err != nil {
		logger.LogError(err, "JSONのデコードに失敗")
		return nil, err
	}

	var images []string
	if err := json.Unmarshal([]byte(row.Images), &images); err != nil {
		logger.LogError(err, "JSONのデコードに失敗")
		return nil, err
	}

	// 日付のパース
	updatedAt, err := time.Parse("2006/01/02", row.UpdateAt)
	if err != nil {
		logger.LogError(err, "日付のパースに失敗")
		return nil, err
	}

	route := &domain.Route{
		ID:               id,
		Title:            row.Title,
		Description:      row.Description,
		FullDescription:  row.FullDescription,
		Distance:         row.Distance,
		Time:             row.Time,
		Tags:             tags,
		Likes:            row.Likes,
		Images:           images,
		Map:              row.Map,
		TotalCheckpoints: row.TotalCheckpoints,
		UpdateAt:         updatedAt,
	}

	// チェックポイントの取得
	checkpoints, err := r.GetCheckpointsByRouteID(id)
	if err != nil {
		logger.LogError(err, "チェックポイントの取得に失敗")
		return nil, err
	}
	route.Checkpoints = checkpoints

	return route, nil
}
