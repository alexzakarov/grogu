package mysql

import (
	"fmt"
	pgConfig "github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/database/mysql"
	postgres2 "github.com/alexzakarov/grogu/database/mysql/base_repo/mysql"
	"github.com/alexzakarov/grogu/database/ports"
	"github.com/alexzakarov/grogu/examples"
	"github.com/alexzakarov/grogu/examples/postgres/config"
	"github.com/alexzakarov/grogu/logger"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func init() {
	var baseDB ports.IBaseDb
	cfg, errConfig := config.ParseConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
	}
	appLogger = logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	tableQuery := `CREATE TABLE IF NOT EXISTS users (
     	user_id int NOT NULL AUTO_INCREMENT,
		meta_data VARCHAR(255) DEFAULT '',
		user_title VARCHAR(75) DEFAULT '',
    	PRIMARY KEY (user_id)
	)`

	// Init Clients
	baseDB, err = mysql.NewMySqlDB(pgConfig.SQLXDbConfig{
		Host:      cfg.Mysql.HOST,
		Port:      cfg.Mysql.PORT,
		User:      cfg.Mysql.USER,
		Pass:      cfg.Mysql.PASS,
		DefaultDb: cfg.Mysql.DEFAULT_DB,
	})
	if err != nil {
		appLogger.Error("Error when tyring to connect to Postgresql")
	} else {
		appLogger.Info("Postgresql connected")
	}

	_, err = baseDB.Exec(tableQuery)
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
	sqlxBaseRepo = postgres2.NewSQLXBaseRepo[examples.CreateUserDbModel, examples.UpdateUserDbModel, examples.UserResDto](pgConfig.PostgresConfig{
		Db: baseDB,
		GeneralConfig: pgConfig.GeneralConfig{
			Table:         "users",
			PrimaryKey:    "user_id",
			SoftDeletable: false,
			StatusName:    "",
			StatusType:    "",
		},
	})
}

func TestSQLXCreate(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	user := examples.CreateUserReqDto{
		UserTitle: "Test User",
	}
	meta := user.ToDbModel("This user has admin role")

	sqlxBaseRepo.Create(meta, func(id int64) {
		record = 1
		userId = id
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})
	fmt.Println(userId)
	assertion.Equal(int64(1), record, "User cannot be created")

}

func TestSQLXUpdate(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	user := examples.UpdateUserReqDto{
		UserTitle: "Test User",
	}
	fmt.Println(userId)
	meta := user.ToDbModel("This user has admin role")
	sqlxBaseRepo.Update(userId, meta, func() {
		record = 1
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(int64(1), record, fmt.Sprintf("User cannot be updated; User ID: %d", userId))

}

func TestSQLXGetOne(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	sqlxBaseRepo.GetOne(userId, func(user examples.UserResDto) {
		record = 1
		_ = user
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(int64(1), record, fmt.Sprintf("failed to retrieve user; User ID: %d", userId))
}

func TestSQLXDeleteOne(t *testing.T) {
	assertion := assert.New(t)

	var (
		record int64
	)

	sqlxBaseRepo.DeleteOne(userId, func() {
		record = 1
	}, func(rec int64) {
		// negative rec refers to db errors
		record = rec
	})

	assertion.Equal(int64(1), record, fmt.Sprintf("User cannot be deleted; User ID: %d", userId))

}

func BenchmarkAllSQLX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestSQLXCreate(&testing.T{})
		TestSQLXUpdate(&testing.T{})
		TestSQLXGetOne(&testing.T{})
		TestSQLXDeleteOne(&testing.T{})
	}
}

func BenchmarkSQLXCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestSQLXCreate(&testing.T{})
	}
}
func BenchmarkSQLXUpdate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestSQLXUpdate(&testing.T{})
	}
}
func BenchmarkSQLXGetOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestSQLXGetOne(&testing.T{})
	}
}

func BenchmarkSQLXDeleteOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestSQLXDeleteOne(&testing.T{})
	}
}
