package config

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
)

type GeneralConfig struct {
	Database      string `json:"database"`
	Schema        string `json:"schema"`
	Table         string `json:"table"`
	PrimaryKey    string `json:"primary_key"`
	SoftDeletable bool   `json:"soft_deletable"`
	StatusName    string `json:"status_name"`
	StatusType    string `json:"status_type"`
}

type PGXBaseRepoConfig struct {
	Ctx context.Context `json:"ctx"`
	Db  *pgxpool.Pool   `json:"db"`
	GeneralConfig
}

type PQBaseRepoConfig struct {
	Db *sql.DB `json:"db"`
	GeneralConfig
}
type SQLXBaseRepoConfig struct {
	Ctx context.Context `json:"ctx"`
	Db  *sqlx.DB        `json:"db"`
	GeneralConfig
}
