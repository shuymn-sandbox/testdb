package handler_test

import (
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/shuymn-sandbox/testdb"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	if os.Getenv("TESTDB_ISOLATE") == "true" {
		db, err = testdb.TestWithMySQL()
	} else {
		config := mysql.Config{
			User:                 os.Getenv("MYSQL_USER"),
			Passwd:               os.Getenv("MYSQL_PASSWORD"),
			Net:                  "tcp",
			Addr:                 os.Getenv("MYSQL_HOST"),
			DBName:               os.Getenv("MYSQL_DATABASE"),
			ParseTime:            true,
			AllowNativePasswords: true,
		}
		db, err = sqlx.Open("mysql", config.FormatDSN())
	}
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}
