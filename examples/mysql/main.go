package mysql

import (
	"context"
	"github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/examples"
	"github.com/alexzakarov/grogu/logger"
)

var (
	ctx = context.Background()

	appLogger    *logger.ApiLogger
	err          error
	sqlxBaseRepo config.IBaseRepo[examples.CreateUserDbModel, examples.UpdateUserDbModel, examples.UserResDto]
	userId       int64
)
