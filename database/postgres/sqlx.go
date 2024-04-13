package postgres

import (
	"fmt"
	"github.com/alexzakarov/grogu/examples/postgres/config"
	"github.com/jmoiron/sqlx"
)

var (
	SQLXConnStr string
)

// NewSQLXPostgresqlDB Return new Postgresql client
func NewSQLXPostgresqlDB(cfg *config.Config) (db *sqlx.DB, err error) {
	println("Driver PostgreSQL Initialized")
	SQLXConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Postgresql.HOST, cfg.Postgresql.PORT, cfg.Postgresql.USER, cfg.Postgresql.PASS, cfg.Postgresql.DEFAULT_DB)

	db, err = sqlx.Connect("postgres", SQLXConnStr)
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
func NewSQLXMySqlDB(cfg *config.Config) (db *sqlx.DB, err error) {
	println("Driver PostgreSQL Initialized")
	SQLXConnStr = fmt.Sprintf("mysql://%s:%s@%s:%d/%s", cfg.Postgresql.USER, cfg.Postgresql.PASS, cfg.Postgresql.HOST, cfg.Postgresql.PORT, cfg.Postgresql.DEFAULT_DB)

	db, err = sqlx.Connect("mysql", SQLXConnStr)
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
