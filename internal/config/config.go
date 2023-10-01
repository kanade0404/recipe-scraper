package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Config struct {
	DB *sql.DB
}

// NewConfig
/*
アプリケーションに必要な設定を定義します。
*/
func NewConfig() (*Config, error) {
	// データベースへのアクセスに必要なオブジェクトを作成する。
	db, err := sql.Open("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" user="+os.Getenv("DB_USER")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &Config{DB: db}, nil
}
