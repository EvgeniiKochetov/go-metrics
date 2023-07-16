package database

import (
	"database/sql"
)

func tableExists(db *sql.DB) error {

	return nil
}

func CreateTable(db *sql.DB) error {
	db.Exec("CREATE TABLE Metrics")
	return nil
}
