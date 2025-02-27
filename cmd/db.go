package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func configDB() *sqlx.DB {
	// sqlx を利用して SQLite データベースに接続
	db, err := sqlx.Open("sqlite3", "./sagaicle.db")
	if err != nil {
		log.Fatalf("SQLite データベースのオープンに失敗しました: %v", err)
	}

	// init.sql の内容をファイルから読み込む
	sqlBytes, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatal(err)
	}
	sqlStatements := string(sqlBytes)

	// 初期化スクリプトを実行
	_, err = db.Exec(sqlStatements)
	if err != nil {
		log.Fatal(err)
	}

	// データベース接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("SQLite データベースへの接続確認に失敗しました: %v", err)
	}

	// 必要なテーブルがなければ作成する
	if err := createTables(db); err != nil {
		log.Fatalf("テーブル作成に失敗しました: %v", err)
	}

	return db
}

// createTables は必要なテーブルを作成します。
// ここでは tags テーブルを例として作成しています。
func createTables(db *sqlx.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS tags (
        tag TEXT PRIMARY KEY,
        CHECK (LENGTH(tag) BETWEEN 1 AND 10)
    );
	`
	_, err := db.Exec(schema)
	return err
}
