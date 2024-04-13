package postgres

import (
	"database/sql"
	"fmt"
	"github.com/alexzakarov/grogu/config"
	_ "github.com/lib/pq"
)

var (
	PQConnStr string
)

// NewPQPostgresqlDB Return new Postgresql client
func NewPQPostgresqlDB(cfg config.PQDbConfig) (db *sql.DB, err error) {
	println("Driver PostgreSQL Initialized")
	PQConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DefaultDb)

	db, err = sql.Open("postgres", PQConnStr)
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
