package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sushigon-dev/sagaicle/utils/logger"
)

func DB() *sqlx.DB {
	// sqlx を利用して SQLite データベースに接続
	db, err := sqlx.Open("sqlite3", "./sagaicle.db")
	if err != nil {
		logger.LogError(err, "DBのオープンに失敗")
		log.Fatalf("SQLite データベースのオープンに失敗しました: %v", err)
	}

	// init.sql の内容をファイルから読み込む
	sqlBytes, err := os.ReadFile("init.sql")
	if err != nil {
		logger.LogError(err, "init.sqlの読み込みに失敗")
		log.Fatal(err)
	}
	sqlStatements := string(sqlBytes)

	// 初期化スクリプトを実行
	_, err = db.Exec(sqlStatements)
	if err != nil {
		logger.LogError(err, "初期化スクリプト実行に失敗")
		log.Fatal(err)
	}

	// データベース接続確認
	if err := db.Ping(); err != nil {
		logger.LogError(err, "DBの接続確認失敗")
		log.Fatalf("SQLite データベースへの接続確認に失敗しました: %v", err)
	}

	return db
}
