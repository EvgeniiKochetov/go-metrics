package config

import (
	"database/sql"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/database"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"sync"
)

var lock = &sync.Mutex{}

type configuration struct {
	db      *sql.DB
	configs map[string]string
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
	if err != nil {
		logger.Log.Info("can't establish connection to database: " + dbConnection)
		return
	}
	c.db = db
	database.DeleteTable(db)
	database.CreateTable(db)
}

func (c *configuration) GetDatabaseConnection() *sql.DB {
	return c.db
}

func (c *configuration) GetFlag(flagName string) string {
	return c.configs[flagName]
}

func (c *configuration) SetFlag(flagName string, flagValue string) error {
	c.configs[flagName] = flagValue
	return nil
}
