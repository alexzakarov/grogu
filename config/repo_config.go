package config

import (
	"context"
	"github.com/alexzakarov/grogu/database/ports"
)

type PostgresConfig struct {
	Ctx context.Context `json:"ctx"`
	Db  ports.IBaseDb   `json:"db"`
	GeneralConfig
}

type GeneralConfig struct {
	Database      string `json:"database"`
	Table         string `json:"table"`
	PrimaryKey    string `json:"primary_key"`
	SoftDeletable bool   `json:"soft_deletable"`
	StatusName    string `json:"status_name"`
	StatusType    string `json:"status_type"`
}

type PGXBaseRepoConfig struct {
	Ctx context.Context `json:"ctx"`
	Db  ports.IBaseDb   `json:"db"`
	GeneralConfig
}

func (c PGXBaseRepoConfig) ToPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Ctx: c.Ctx,
		Db:  c.Db,
		GeneralConfig: GeneralConfig{
			Database:      c.Database,
			Table:         c.Table,
			PrimaryKey:    c.PrimaryKey,
			SoftDeletable: c.SoftDeletable,
			StatusName:    c.StatusName,
			StatusType:    c.StatusType,
		},
	}
}

type PQBaseRepoConfig struct {
	Db ports.IBaseDb `json:"db"`
	GeneralConfig
}

func (c PQBaseRepoConfig) ToPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Db: c.Db,
		GeneralConfig: GeneralConfig{
			Database:      c.Database,
			Table:         c.Table,
			PrimaryKey:    c.PrimaryKey,
			SoftDeletable: c.SoftDeletable,
			StatusName:    c.StatusName,
			StatusType:    c.StatusType,
		},
	}
}

type SQLXBaseRepoConfig struct {
	Db ports.IBaseDb `json:"db"`
	GeneralConfig
}

func (c SQLXBaseRepoConfig) ToPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Db: c.Db,
		GeneralConfig: GeneralConfig{
			Database:      c.Database,
			Table:         c.Table,
			PrimaryKey:    c.PrimaryKey,
			SoftDeletable: c.SoftDeletable,
			StatusName:    c.StatusName,
			StatusType:    c.StatusType,
		},
	}
}
