package config

import (
	"database/sql"
	"fmt"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"
	_ "github.com/jackc/pgx"
	"sync"
)

var lock = &sync.Mutex{}

type configuration struct {
	db *sql.DB
}

var instanceOfConf *configuration = nil

func GetInstance() *configuration {
	if instanceOfConf == nil {
		lock.Lock()
		defer lock.Unlock()

		if instanceOfConf == nil {
			instanceOfConf = &configuration{}
		}
	}

	return instanceOfConf
}

func (c *configuration) SetDB(dbConnection string) {
	db, err := sql.Open("pgx", dbConnection)
	fmt.Println("setDB:", db, "/string:", dbConnection)
	if err != nil {
		fmt.Println(err)
		logger.Log.Info("can't establish connection to database: " + dbConnection)
		return
	}
	c.db = db
	defer db.Close()
}

func (c *configuration) GetDatabaseConnection() *sql.DB {
	return c.db
}
