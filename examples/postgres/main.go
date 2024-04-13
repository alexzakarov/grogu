package postgres

import (
	"context"
	"database/sql"
	repo "github.com/alexzakarov/grogu/base_repo/postgres"
	"github.com/alexzakarov/grogu/examples"
	"github.com/alexzakarov/grogu/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
)

var (
	ctx = context.Background()

	pgxDb        *pgxpool.Pool
	pqDb         *sql.DB
	sqlxDb       *sqlx.DB
	appLogger    *logger.ApiLogger
	err          error
	pgxBaseRepo  repo.IBaseRepo[examples.CreateUserDbModel, examples.UpdateUserDbModel, examples.UserResDto]
	pqBaseRepo   repo.IBaseRepo[examples.CreateUserDbModel, examples.UpdateUserDbModel, examples.UserResDto]
	sqlxBaseRepo repo.IBaseRepo[examples.CreateUserDbModel, examples.UpdateUserDbModel, examples.UserResDto]
	userId       int64
)
