package mysql

import (
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
	var cmd sql.Result
	cmd, err = s.db.Exec(query, dat...)
	if err != nil {
		return 0, err
	}
	data, err = cmd.LastInsertId()
	fmt.Println("test: ", data)
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

// NewMySqlDB Return new Postgresql client
func NewMySqlDB(cfg config.SQLXDbConfig) (db ports.IBaseDb, err error) {
	var baseDB *sqlx.DB
	println("Driver mysql Initialized")
	SQLXConnStr = fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DefaultDb)

	baseDB, err = sqlx.Connect("mysql", SQLXConnStr)
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
