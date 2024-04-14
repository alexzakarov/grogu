package postgres

import (
	"database/sql"
	"fmt"
	"github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/database/ports"
	_ "github.com/lib/pq"
)

var (
	PQConnStr string
)

type SqlDb struct {
	db *sql.DB
}

func (s *SqlDb) Exec(query string, dat ...interface{}) (rows_affected int64, err error) {
	var cmd sql.Result

	cmd, err = s.db.Exec(query, dat...)
	if err != nil {
		return 0, err
	}
	rows_affected, err = cmd.RowsAffected()
	if err != nil {
		return 0, err
	}

	return
}

func (s *SqlDb) Insert(query string, dat ...interface{}) (data int64, err error) {

	err = s.db.QueryRow(query, dat...).Scan(&data)

	return
}

func (s *SqlDb) Update(query string, dat ...interface{}) (rows_affected int64, err error) {

	var cmd sql.Result

	cmd, err = s.db.Exec(query, dat...)
	if err != nil {
		return 0, err
	}
	rows_affected, err = cmd.RowsAffected()
	if err != nil {
		return 0, err
	}

	return
}

func (s *SqlDb) Select(query string, dat ...interface{}) (data []byte, err error) {

	err = s.db.QueryRow(query, dat...).Scan(&data)

	return
}

func (s *SqlDb) Delete(query string, dat ...interface{}) (rows_affected int64, err error) {

	var cmd sql.Result

	cmd, err = s.db.Exec(query, dat...)
	if err != nil {
		return 0, err
	}
	rows_affected, err = cmd.RowsAffected()
	if err != nil {
		return 0, err
	}

	return
}

// NewPQPostgresqlDB Return new Postgresql client
func NewPQPostgresqlDB(cfg config.PQDbConfig) (db ports.IBaseDb, err error) {
	var baseDB *sql.DB
	println("Driver PostgreSQL Initialized")
	PQConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DefaultDb)

	baseDB, err = sql.Open("postgres", PQConnStr)
	if err != nil {
		println(err.Error())
	} else {
		print("conn ok")
	}

	if err = baseDB.Ping(); err != nil {
		println(err.Error())
	}
	return &SqlDb{
		db: baseDB,
	}, err
}
