package config

import (
	"database/sql"
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

func (c *configuration) CheckConnection() error {
	if err := c.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *configuration) SetDB(db *sql.DB) {
	c.db = db
}
