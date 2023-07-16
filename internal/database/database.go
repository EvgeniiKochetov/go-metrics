package database

import (
	"database/sql"
	"fmt"
)

func tableExists(db *sql.DB) (bool, error) {
	r := "SELECT EXISTS (SELECT * FROM pg_tables WHERE tablename = 'Metrics')"
	_, err := db.Exec(r)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateTable(db *sql.DB) error {
	createTable := "DELETE TABLE Metics; CREATE TABLE Metrics (type varchar(10), name varchar(50), value DOUBLE PRECISION, counter integer)"

	_, err := db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func AddGaugeMetric(db *sql.DB, nameOfMetric string, value float64) error {
	return nil
}

func AddCounterMetric(db *sql.DB, nameOfMetric string, value float64) error {
	return nil
}

func GetGaugeMetric(db *sql.DB, nameOfMetric string) (*float64, error) {
	pointer := float64(0)
	return &pointer, nil
}

func GetCounterMetric(db *sql.DB, nameOfMetric string) (*int64, error) {
	pointer := int64(0)
	return &pointer, nil
}
