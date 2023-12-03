package testdb

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

//go:embed example/schema.sql
var schema string

func TestWithMySQL() (*sqlx.DB, error) {
	// 適当にテスト用のランダムなDB名を生成
	dbName := fmt.Sprintf("%s_%s_test", os.Getenv("MYSQL_DATABASE"), generateRandomString())
	if err := createMySQLDatabase(dbName); err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	config := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               dbName,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db := sqlx.MustOpen("mysql", config.FormatDSN())
	// テスト用DBへのマイグレーション
	for _, query := range strings.Split(schema, ";") {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if _, err := db.DB.Exec(query); err != nil {
			return nil, fmt.Errorf("failed to migrate: %w", err)
		}
	}
	return db, nil
}

func createMySQLDatabase(dbName string) (err error) {
	// 権限の問題があるので、rootユーザーで接続する
	config := mysql.Config{
		User:   "root",
		Passwd: os.Getenv("MYSQL_ROOT_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("MYSQL_HOST"),
		// とりあえずデフォルトのテーブル名を使う
		DBName:               os.Getenv("MYSQL_DATABASE"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db := sqlx.MustOpen("mysql", config.FormatDSN())
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = multierr.Append(err, fmt.Errorf("failed to close db: %w", closeErr))
		}
	}()

	// テスト用のDBを作成
	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	// テストで使うユーザーに権限を付与
	if _, err = db.Exec(fmt.Sprintf("GRANT ALL ON %s.* TO %s@'%%'", dbName, os.Getenv("MYSQL_USER"))); err != nil {
		return fmt.Errorf("failed to grant privileges: %w", err)
	}
	return nil
}
