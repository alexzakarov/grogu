package postgres

import (
	"context"
	"fmt"
	"github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/database/ports"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ConnStr string
)

const (
	maxOpenConns    = 250
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

type PgxDb struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func (s *PgxDb) Exec(query string, dat ...interface{}) (rows_affected int64, err error) {
	var cmd pgconn.CommandTag

	cmd, err = s.db.Exec(s.ctx, query, dat...)
	if err != nil {
		return 0, err
	}
	rows_affected = cmd.RowsAffected()
	return
}

func (s *PgxDb) Insert(query string, dat ...interface{}) (data int64, err error) {

	err = s.db.QueryRow(s.ctx, query, dat...).Scan(&data)

	return
}

func (s *PgxDb) Update(query string, dat ...interface{}) (rows_affected int64, err error) {

	var cmd pgconn.CommandTag

	cmd, err = s.db.Exec(s.ctx, query, dat...)
	if err != nil {
		return 0, err
	}
	rows_affected = cmd.RowsAffected()

	return
}

func (s *PgxDb) Select(query string, dat ...interface{}) (data []byte, err error) {

	err = s.db.QueryRow(s.ctx, query, dat...).Scan(&data)

	return
}

func (s *PgxDb) Delete(query string, dat ...interface{}) (rows_affected int64, err error) {

	var cmd pgconn.CommandTag

	cmd, err = s.db.Exec(s.ctx, query, dat...)
	if err != nil {
		return 0, err
	}
	rows_affected = cmd.RowsAffected()

	return
}

// NewPGXPostgresqlDB Return new Postgresql client
func NewPGXPostgresqlDB(cfg config.PGXDbConfig) (db ports.IBaseDb, err error) {
	var baseDB *pgxpool.Pool
	ctx := context.Background()
	println("Driver PostgreSQL Initialized")
	ConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DefaultDb, cfg.MaxConn)

	baseDB, err = pgxpool.Connect(ctx, ConnStr)
	if err != nil {
		println(err.Error())
	} else {
		print("conn ok")
	}

	if err = baseDB.Ping(context.Background()); err != nil {
		println(err.Error())
	}
	return &PgxDb{
		ctx: ctx,
		db:  baseDB,
	}, err
}
