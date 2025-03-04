package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
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
		logger.Error(err, "ルートの登録に失敗"+fmt.Sprint(query,
			route.ID.String(), route.Title, route.Description, route.FullDescription,
			route.Distance, route.Time, route.Likes, image, updateAt,
			route.TotalCheckpoints, route.Map),
		)
		return err
	}

	for _, tag := range route.Tags {
		// タグマスターへの登録：tags テーブルにタグを追加（既に存在する場合は無視）
		tagInsertQuery := `INSERT OR IGNORE INTO tags (tag) VALUES (?);`
		if _, err := r.db.Exec(tagInsertQuery, tag); err != nil {
			logger.Error(err, "タグマスターへの登録に失敗: "+tag)
			return err
		}

		// ルートとタグの紐付け：Route_Tagsテーブルに対して、ルートIDとタグ名を登録
		routeTagQuery := `INSERT INTO route_tags (route_id, tag_name) VALUES (?, ?);`
		if _, err = r.db.Exec(routeTagQuery, route.ID.String(), tag); err != nil {
			logger.Error(err, "ルートタグの登録に失敗: "+
				fmt.Sprint(routeTagQuery, route.ID.String(), tag))
			return err
		}
	}

	// 画像の登録：Route_Imagesテーブルに対して、ルートIDと画像URIを登録する
	for _, imageURI := range route.Images {
		imageQuery := `INSERT INTO route_images (route_id, image) VALUES (?, ?);`
		if _, err = r.db.Exec(imageQuery, route.ID.String(), imageURI); err != nil {
			logger.Error(err, "画像の登録に失敗: "+imageURI)
			return err
		}
	}

	// チェックポイントの登録
	return r.CreateCheckpoints(route.ID, route.Checkpoints)
}

// 指定したルート ID の詳細情報を取得
func (r *SQLiteRepository) GetRouteByID(id uuid.UUID) (*domain.Route, error) {
	// routes テーブルから取得するクエリ
	query := `
        SELECT id, title, description, full_description, distance, time, likes, image, update_at, total_checkpoints, map
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
		Likes            int     `db:"likes"`
		Image            string  `db:"image"`
		UpdateAt         string  `db:"update_at"`
		TotalCheckpoints int     `db:"total_checkpoints"`
		Map              string  `db:"map"`
	}

	if err := r.db.Get(&row, query, id.String()); err != nil {
		logger.Error(err, "ルートの取得に失敗"+fmt.Sprint(query, id.String()))
		return nil, err
	}

	// 日付のパース
	updatedAt, err := time.Parse("2006/01/02", row.UpdateAt)
	if err != nil {
		logger.Error(err, "日付のパースに失敗")
		return nil, err
	}

	route := &domain.Route{
		ID:              id,
		Title:           row.Title,
		Description:     row.Description,
		FullDescription: row.FullDescription,
		Distance:        row.Distance,
		Time:            row.Time,
		Likes:           row.Likes,
		// ルートサマリー用の image は routes テーブルの image を利用
		Images:           nil, // 後で取得
		Map:              row.Map,
		TotalCheckpoints: row.TotalCheckpoints,
		UpdateAt:         updatedAt,
	}

	// タグの取得（route_tags テーブルから）
	tagsQuery := `SELECT tag_name FROM route_tags WHERE route_id = ?;`
	var tags []string
	if err := r.db.Select(&tags, tagsQuery, id.String()); err != nil {
		logger.Error(err, "タグの取得に失敗")
		return nil, err
	}
	route.Tags = tags

	// 画像の取得（route_images テーブルから）
	imagesQuery := `SELECT image FROM route_images WHERE route_id = ?;`
	var images []string
	if err := r.db.Select(&images, imagesQuery, id.String()); err != nil {
		logger.Error(err, "画像の取得に失敗")
		return nil, err
	}
	route.Images = images

	// チェックポイントの取得
	checkpoints, err := r.GetCheckpointsByRouteID(id)
	if err != nil {
		logger.Error(err, "チェックポイントの取得に失敗")
		return nil, err
	}
	route.Checkpoints = checkpoints

	return route, nil
}

// 検索条件に基づいてルート一覧とヒット件数を返す（検索する）
// リクエストはポインタで受け取ることで返り値を略
func (r *SQLiteRepository) SearchRoutes(criteria *domain.SearchCriteria) ([]domain.RouteSummary, int, error) {
	var whereClauses []string
	var args []interface{}

	// 距離条件: 条件として指定しない場合は -1.0 を用いる（初期値）
	if criteria.Distance.Min != -1.0 {
		whereClauses = append(whereClauses, "distance >= ?")
		args = append(args, criteria.Distance.Min)
	}
	if criteria.Distance.Max != -1.0 {
		whereClauses = append(whereClauses, "distance <= ?")
		args = append(args, criteria.Distance.Max)
	}

	// 所要時間条件: 条件として指定しない場合は -1 を用いる（初期値）
	if criteria.Time.Min != -1 {
		whereClauses = append(whereClauses, "time >= ?")
		args = append(args, criteria.Time.Min)
	}
	if criteria.Time.Max != -1 {
		whereClauses = append(whereClauses, "time <= ?")
		args = append(args, criteria.Time.Max)
	}

	// タグ条件: タグが指定されている場合、検索オプションに応じたサブクエリでフィルタ
	if len(criteria.Tags) > 0 {
		// 検索オプションのデフォルトは "OR"
		searchOption := "OR"
		if len(criteria.SearchOption) > 0 {
			searchOption = criteria.SearchOption[0]
		}

		// タグリスト用のプレースホルダーを作成
		placeholders := make([]string, len(criteria.Tags))
		for i := range criteria.Tags {
			placeholders[i] = "?"
		}

		// サブクエリの組み立て
		tagListClause := strings.Join(placeholders, ",")
		switch searchOption {
		case "OR":
			// ルートが指定されたタグのいずれかを持っていればマッチ
			whereClauses = append(whereClauses, "id IN (SELECT route_id FROM route_tags WHERE tag_name IN ("+tagListClause+"))")
			for _, tag := range criteria.Tags {
				args = append(args, tag)
			}
		case "AND":

			// ルートが指定された全てのタグを持っている必要がある
			whereClauses = append(whereClauses, "id IN (SELECT route_id FROM route_tags WHERE tag_name IN ("+tagListClause+") GROUP BY route_id HAVING COUNT(DISTINCT tag_name) = ?)")
			for _, tag := range criteria.Tags {
				args = append(args, tag)
			}
			args = append(args, len(criteria.Tags))
		case "NOT":

			// ルートが指定されたタグのいずれも持たない
			whereClauses = append(whereClauses, "id NOT IN (SELECT route_id FROM route_tags WHERE tag_name IN ("+tagListClause+"))")
			for _, tag := range criteria.Tags {
				args = append(args, tag)
			}
		}
	}

	// WHERE 句の組み立て: 条件がない場合は常に真となる "1=1" を利用
	whereSQL := "1=1"
	if len(whereClauses) > 0 {
		whereSQL = strings.Join(whereClauses, " AND ")
	}

	// ヒット件数を取得するクエリ
	countQuery := "SELECT COUNT(*) FROM routes WHERE " + whereSQL
	var hitCount int
	if err := r.db.Get(&hitCount, countQuery, args...); err != nil {
		logger.Error(err, "ヒット件数の取得に失敗"+fmt.Sprint(countQuery, args))
		return nil, 0, err
	}

	// メインクエリの組み立て:
	// 各ルートのタグは、サブクエリで GROUP_CONCAT を使って連結し、Go 側で分割
	mainQuery := "SELECT id, title, description, distance, time, likes, image, update_at, " +
		"(SELECT GROUP_CONCAT(tag_name, '|||') FROM route_tags WHERE route_id = routes.id) as tags " +
		"FROM routes WHERE " + whereSQL

	// ソートキー
	sortKey := criteria.Sort.Key
	switch sortKey {
	case "distance", "time", "likes", "update_at":
		// 有効なキーの場合はそのまま利用
	default:
		sortKey = "likes"
	}

	// ソート順
	sortOrder := criteria.Sort.Order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	mainQuery += " ORDER BY " + sortKey + " " + sortOrder

	// 件数制限
	if criteria.Limit < 1 {
		criteria.Limit = 12
	} else if criteria.Limit > 60 {
		criteria.Limit = 60
	}
	mainQuery += " LIMIT ?"
	args = append(args, criteria.Limit)

	// クエリ実行: 中間構造体を利用して結果を一時受け取り
	var rows []struct {
		ID          string         `db:"id"`
		Title       string         `db:"title"`
		Description string         `db:"description"`
		Distance    float64        `db:"distance"`
		Time        int            `db:"time"`
		Likes       int            `db:"likes"`
		Image       string         `db:"image"`
		UpdateAt    string         `db:"update_at"`
		Tags        sql.NullString `db:"tags"`
	}
	if err := r.db.Select(&rows, mainQuery, args...); err != nil {
		return nil, 0, err
	}

	// 結果のマッピング
	var summaries []domain.RouteSummary
	for _, row := range rows {
		var tags []string

		// タグが存在する場合は分割
		if row.Tags.Valid && row.Tags.String != "" {
			tags = strings.Split(row.Tags.String, "|||")
		}

		// JSON 文字列をTime型に変換
		updatedAt, err := time.Parse("2006/01/02", row.UpdateAt)
		if err != nil {
			logger.Error(err, "日付のパースに失敗")
			return nil, 0, err
		}

		// サマリー用に画像は配列の最初の要素を利用
		summary := domain.RouteSummary{
			ID:          uuid.MustParse(row.ID),
			Title:       row.Title,
			Description: row.Description,
			Distance:    row.Distance,
			Time:        row.Time,
			Tags:        tags,
			Likes:       row.Likes,
			Image:       row.Image,
			UpdateAt:    updatedAt, // 既に "yyyy/mm/dd" 形式で保存されている前提
		}
		summaries = append(summaries, summary)
	}

	return summaries, hitCount, nil
}
