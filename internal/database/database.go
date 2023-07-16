package database

import (
	"database/sql"
	"fmt"
)

func tableExists(db *sql.DB) error {

	return nil
}

func CreateTable(db *sql.DB) error {
	createTable := fmt.Sprintf("CREATE TABLE Metrics (type varchar(10), name varchar(50), value DOUBLE PRECISION, counter integer)")

	_, err := db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
