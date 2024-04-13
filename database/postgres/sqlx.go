package postgres

import (
	"context"
	"fmt"
	"github.com/alexzakarov/grogu/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	SQLXConnStr string
)

// NewSQLXPostgresqlDB Return new Postgresql client
func NewSQLXPostgresqlDB(cfg config.SQLXDbConfig) (db *sqlx.DB, err error) {
	println("Driver PostgreSQL Initialized")
	SQLXConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DefaultDb)

	db, err = sqlx.ConnectContext(context.Background(), "postgres", SQLXConnStr)
	if err != nil {
		println(err.Error())
	} else {
		print("conn ok")
	}

	if err = db.Ping(); err != nil {
		println(err.Error())
	}
	return
}

// NewSQLXPostgresqlDB Return new Postgresql client
func NewSQLXMySqlDB(cfg config.SQLXDbConfig) (db *sqlx.DB, err error) {
	println("Driver PostgreSQL Initialized")
	SQLXConnStr = fmt.Sprintf("mysql://%s:%s@%s:%d/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DefaultDb)

	db, err = sqlx.ConnectContext(context.Background(), "mysql", SQLXConnStr)
	if err != nil {
		println(err.Error())
	} else {
		print("conn ok")
	}

	if err = db.Ping(); err != nil {
		println(err.Error())
	}
	return
}
