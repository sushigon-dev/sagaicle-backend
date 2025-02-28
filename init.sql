-- タグマスタテーブル
CREATE TABLE IF NOT EXISTS tags (
    tag TEXT PRIMARY KEY,
    CHECK (LENGTH(tag) BETWEEN 1 AND 10)
);

-- タグ数上限(20件)を enforce するトリガー
CREATE TRIGGER IF NOT EXISTS limit_tags_count
BEFORE INSERT ON tags
WHEN (SELECT COUNT(*) FROM tags) >= 20
BEGIN
    SELECT RAISE(FAIL, 'Maximum number of tags (20) reached.');
END;

-- ルートテーブル
CREATE TABLE IF NOT EXISTS routes (
    id TEXT PRIMARY KEY,  -- UUIDv4 を想定
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    full_description TEXT NOT NULL,
    distance REAL NOT NULL,
    time INTEGER NOT NULL,
    likes INTEGER NOT NULL DEFAULT 0,
    image TEXT NOT NULL,
    update_at TEXT NOT NULL,  -- yyyy/mm/dd の文字列
    total_checkpoints INTEGER NOT NULL,
    map TEXT NOT NULL,
    CHECK (LENGTH(title) BETWEEN 1 AND 20),
    CHECK (LENGTH(description) BETWEEN 1 AND 60),
    CHECK (LENGTH(full_description) BETWEEN 1 AND 200),
    CHECK (distance >= 0),
    CHECK (time >= 0),
    CHECK (total_checkpoints BETWEEN 1 AND 20),
    CHECK (LENGTH(image) BETWEEN 8 AND 1023),
    CHECK (image LIKE 'https://%'),
    CHECK (LENGTH(update_at) = 10),
    CHECK (update_at GLOB '[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]'),
    CHECK (LENGTH(map) BETWEEN 8 AND 1023),
    CHECK (map LIKE 'https://www.google.com%')
);

-- ルートの画像一覧を管理するテーブル
CREATE TABLE IF NOT EXISTS route_images (
    route_id TEXT NOT NULL,
    image TEXT NOT NULL,
    PRIMARY KEY (route_id, image),
    FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE,
    CHECK (LENGTH(image) BETWEEN 8 AND 1023),
    CHECK (image LIKE 'https://%')
);

-- ルート1件につき最大6件の画像登録を enforce するトリガー
CREATE TRIGGER IF NOT EXISTS limit_route_images_insert
BEFORE INSERT ON route_images
FOR EACH ROW
WHEN ((SELECT COUNT(*) FROM route_images WHERE route_id = NEW.route_id) >= 6)
BEGIN
    SELECT RAISE(FAIL, 'Maximum number of images (6) reached for this route.');
END;

-- ルートに付与するタグの中間テーブル
CREATE TABLE IF NOT EXISTS route_tags (
    route_id TEXT NOT NULL,
    tag_name TEXT NOT NULL,
    PRIMARY KEY (route_id, tag_name),
    FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_name) REFERENCES tags(tag) ON DELETE CASCADE,
    CHECK (LENGTH(tag_name) BETWEEN 1 AND 10)
);

-- ルート1件につき最大20件のタグ登録を enforce するトリガー
CREATE TRIGGER IF NOT EXISTS limit_route_tags_insert
BEFORE INSERT ON route_tags
FOR EACH ROW
WHEN ((SELECT COUNT(*) FROM route_tags WHERE route_id = NEW.route_id) >= 20)
BEGIN
    SELECT RAISE(FAIL, 'Maximum number of tags (20) reached for this route.');
END;

-- チェックポイントテーブル（各ルートに対するチェックポイント）
CREATE TABLE IF NOT EXISTS checkpoints (
    route_id TEXT NOT NULL,
    checkpoint_index INTEGER NOT NULL,  -- 順序を表す（例: 0～total_checkpoints-1）
    name TEXT NOT NULL,
    lat REAL NOT NULL,
    lng REAL NOT NULL,
    PRIMARY KEY (route_id, checkpoint_index),
    FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE,
    CHECK (checkpoint_index >= 0)
);

-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,  -- UUIDv4 を想定
    user_name TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    mileage REAL NOT NULL DEFAULT 0,
    total_distance REAL NOT NULL,
    CHECK (LENGTH(user_name) BETWEEN 1 AND 32)
);

-- いいね（ルートに対するユーザーのいいね）中間テーブル
CREATE TABLE IF NOT EXISTS route_likes (
    user_id TEXT NOT NULL,
    route_id TEXT NOT NULL,
    PRIMARY KEY (user_id, route_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE
);

-- バッジ（ルートに対して取得したバッジ）中間テーブル
CREATE TABLE IF NOT EXISTS route_badges (
    user_id TEXT NOT NULL,
    route_id TEXT NOT NULL,
    PRIMARY KEY (user_id, route_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE
);

-- 訪問済みチェックポイント（ユーザーが訪問したチェックポイント）中間テーブル
CREATE TABLE IF NOT EXISTS visited_checkpoints (
    user_id TEXT NOT NULL,
    route_id TEXT NOT NULL,
    checkpoint_index INTEGER NOT NULL,
    PRIMARY KEY (user_id, route_id, checkpoint_index),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (route_id, checkpoint_index) REFERENCES checkpoints(route_id, checkpoint_index) ON DELETE CASCADE
);

