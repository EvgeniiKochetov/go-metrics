package database

import (
	"database/sql"
	"fmt"
)

func tableExists(db *sql.DB) error {

	return nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE Metrics;")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
