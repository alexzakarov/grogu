package postgres

import (
	"database/sql"
	"fmt"
	"github.com/alexzakarov/grogu/examples/postgres/config"
	_ "github.com/lib/pq"
)

var (
	PQConnStr string
)

// NewPQPostgresqlDB Return new Postgresql client
func NewPQPostgresqlDB(cfg *config.Config, driver_name string) (db *sql.DB, err error) {
	println("Driver PostgreSQL Initialized")
	PQConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Postgresql.HOST, cfg.Postgresql.PORT, cfg.Postgresql.USER, cfg.Postgresql.PASS, cfg.Postgresql.DEFAULT_DB)

	db, err = sql.Open(driver_name, PQConnStr)
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
