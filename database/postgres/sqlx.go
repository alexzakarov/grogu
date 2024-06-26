package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/database/ports"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	SQLXConnStr string
)

type SqlxDb struct {
	db *sqlx.DB
}

func (s *SqlxDb) Exec(query string, dat ...interface{}) (rows_affected int64, err error) {
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

func (s *SqlxDb) Insert(query string, dat ...interface{}) (data int64, err error) {

	err = s.db.QueryRow(query, dat...).Scan(&data)

	return
}

func (s *SqlxDb) Update(query string, dat ...interface{}) (rows_affected int64, err error) {

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

func (s *SqlxDb) Select(query string, dat ...interface{}) (data []byte, err error) {

	err = s.db.QueryRow(query, dat...).Scan(&data)

	return
}

func (s *SqlxDb) Delete(query string, dat ...interface{}) (rows_affected int64, err error) {

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

// NewSQLXPostgresqlDB Return new Postgresql client
func NewSQLXPostgresqlDB(cfg config.SQLXDbConfig) (db ports.IBaseDb, err error) {
	var baseDB *sqlx.DB
	println("Driver PostgreSQL Initialized")
	SQLXConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DefaultDb)

	baseDB, err = sqlx.ConnectContext(context.Background(), "postgres", SQLXConnStr)
	if err != nil {
		println(err.Error())
	} else {
		print("conn ok")
	}

	if err = baseDB.Ping(); err != nil {
		println(err.Error())
	}
	return &SqlxDb{
		db: baseDB,
	}, err
}
