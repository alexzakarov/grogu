package postgres

import (
	"fmt"
	postgres2 "github.com/alexzakarov/grogu/base_repo/postgres"
	pgConfig "github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/database/postgres"
	"github.com/alexzakarov/grogu/examples"
	"github.com/alexzakarov/grogu/examples/postgres/config"
	"github.com/alexzakarov/grogu/logger"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
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
		user_title VARCHAR(75) DEFAULT NULL
	)`

	// Init Clients
	pqDb, err = postgres.NewPQPostgresqlDB(cfg, "postgres")
	if err != nil {
		appLogger.Error("Error when tyring to connect to Postgresql")
	} else {
		appLogger.Info("Postgresql connected")
	}

	_, err = pqDb.Exec(tableQuery)
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
	pqBaseRepo = postgres2.NewPQBaseRepo[examples.CreateUserDbModel, examples.UpdateUserDbModel, examples.UserResDto](pgConfig.PQBaseRepoConfig{
		Db: pqDb,
		GeneralConfig: pgConfig.GeneralConfig{
			Schema:        "public",
			Table:         "users",
			PrimaryKey:    "user_id",
			SoftDeletable: false,
			StatusName:    "",
			StatusType:    "",
		},
	})
}

func TestPQCreate(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	user := examples.CreateUserReqDto{
		UserTitle: "Test User",
	}
	meta := user.ToDbModel("This user has admin role")

	pqBaseRepo.Create(meta, func(id int64) {
		record = 1
		userId = id
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(int64(1), record, "User cannot be created")

}

func TestPQUpdate(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	user := examples.UpdateUserReqDto{
		UserTitle: "Test User",
	}
	meta := user.ToDbModel("This user has admin role updated")
	pqBaseRepo.Update(userId, meta, func() {
		record = 1
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(int64(1), record, fmt.Sprintf("User cannot be updated; User ID: %d", userId))

}

func TestPQGetOne(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	pqBaseRepo.GetOne(userId, func(user examples.UserResDto) {
		record = 1
		_ = user
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(int64(1), record, fmt.Sprintf("failed to retrieve user; User ID: %d", userId))
}

func TestPQDeleteOne(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	pqBaseRepo.DeleteOne(userId, func() {
		record = 1
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(int64(1), record, fmt.Sprintf("User cannot be deleted; User ID: %d", userId))

}

func BenchmarkAllPQ(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestPQCreate(&testing.T{})
		TestPQUpdate(&testing.T{})
		TestPQGetOne(&testing.T{})
		TestPQDeleteOne(&testing.T{})
	}
}

func BenchmarkPQCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestPQCreate(&testing.T{})
	}
}
func BenchmarkPQUpdate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestPQUpdate(&testing.T{})
	}
}
func BenchmarkPQGetOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestPQGetOne(&testing.T{})
	}
}

func BenchmarkPQDeleteOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestPQDeleteOne(&testing.T{})
	}
}