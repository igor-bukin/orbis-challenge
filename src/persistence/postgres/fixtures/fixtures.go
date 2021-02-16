package fixtures

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq" // import postgres driver for sql.Open()
	"github.com/orbis-challenge/src/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/testfixtures.v2"
)

var (
	p  *pg.DB
	db *sql.DB

	fixtures *testfixtures.Context
)

func init() { // nolint
	err := config.Load("../../../config.json")
	if err != nil {
		logrus.Fatal("failed to load config", "error", err)
	}

	dbConf := config.Config.PostgresTest

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName))
	if err != nil {
		logrus.Fatal("failed to open sql.DBName with test db config", "error", err)
	}

	// getting location of fixtures path
	// nolint
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename)
	// creating the context that hold the fixtures
	fixtures, err = testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, filepath.Join(path, "data"))
	if err != nil {
		logrus.Fatal("failed to load postgres test db fixtures", "error", err)
	}

	p = pg.Connect(&pg.Options{
		User:     dbConf.User,
		Password: dbConf.Password,
		Database: dbConf.DBName,
		Addr:     dbConf.Host + ":" + dbConf.Port,
	})
}

func GetDB() *pg.DB {
	return p
}

func PrepareFixtures() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}
