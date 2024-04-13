package postgres

import (
	"context"
	"fmt"
	"github.com/alexzakarov/grogu/config"
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

// NewPGXPostgresqlDB Return new Postgresql client
func NewPGXPostgresqlDB(cfg config.PGXDbConfig) (db *pgxpool.Pool, err error) {
	println("Driver PostgreSQL Initialized")
	ConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DefaultDb, cfg.MaxConn)

	db, err = pgxpool.Connect(context.Background(), ConnStr)
	if err != nil {
		println(err.Error())
	} else {
		print("conn ok")
	}

	if err = db.Ping(context.Background()); err != nil {
		println(err.Error())
	}

	return
}
