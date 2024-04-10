package examples

import (
	"context"
	"fmt"
	"github.com/alexzakarov/grogu/base_repo"
	"github.com/alexzakarov/grogu/database"
	"github.com/alexzakarov/grogu/examples/config"
	"github.com/alexzakarov/grogu/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	ctx       = context.Background()
	db        *pgxpool.Pool
	appLogger *logger.ApiLogger
	err       error
	repo      base_repo.IBaseRepo[CreateUserDbModel, UpdateUserDbModel, UserResDto]
)

func init() {
	cfg, errConfig := config.ParseConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
	}
	appLogger = logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	tableQuery := `CREATE TABLE IF NOT EXISTS public.users (
     	user_id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		meta_data VARCHAR(255) DEFAULT NULL,
		user_title VARCHAR(75) DEFAULT NULL,
	)`

	// Init Clients
	db, err = database.NewPostgresqlDB(cfg)
	if err != nil {
		appLogger.Error("Error when tyring to connect to Postgresql")
	} else {
		appLogger.Info("Postgresql connected")
	}

	_, err = db.Exec(ctx, tableQuery)
	if err != nil {
		println(err.Error())
		return
	}

	//status_type gets two values which can be "int" or "bool"
	//status int:
	//    1 - Active
	//    2 - Passive
	//    3 - Block
	//    4 - Delete
	//status bool:
	//    1 - Active
	//    2 - Passive
	repo = base_repo.NewBaseRepo[CreateUserDbModel, UpdateUserDbModel, UserResDto](ctx, db, "public", "users", "user_id", false, "", "")
}

func TestCreate(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	user := CreateUserReqDto{
		UserTitle: "Test User",
	}
	meta := user.ToDbModel("This user has admin role")

	repo.Create(meta, func(id int64) {
		record = 1
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(record, 1, "User cannot be created")

}

func TestUpdate(t *testing.T) {
	assertion := assert.New(t)

	var (
		userId = int64(1)
		record int64
	)

	user := UpdateUserReqDto{
		UserTitle: "Test User",
	}
	meta := user.ToDbModel("This user has admin role")
	repo.Update(userId, meta, func() {
		record = 1
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(record, 1, "User cannot be updated")

}

func TestGetOne(t *testing.T) {
	assertion := assert.New(t)

	var (
		resData UserResDto
		userId  = int64(1)
		record  int64
	)

	repo.GetOne(userId, func(user UserResDto) {
		record = 1
		resData = user
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(record, 1, "User cannot be got")
	fmt.Println(resData)
}

func TestDeleteOne(t *testing.T) {
	assertion := assert.New(t)

	var (
		userId = int64(1)
		record int64
	)

	repo.DeleteOne(userId, func() {
		record = 1
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(record, 1, "User cannot be got")

}