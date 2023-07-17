package database

import (
	"database/sql"
	"fmt"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"
	"strconv"
)

func tableExists(db *sql.DB) (bool, error) {
	r := "SELECT EXISTS (SELECT * FROM pg_tables WHERE tablename = 'Metrics')"
	_, err := db.Exec(r)
	if err != nil {
		return false, err
	}
	return true, nil
}
func DeleteTable(db *sql.DB) error {
	deleteTable := "DROP TABLE Metrics"

	_, err := db.Exec(deleteTable)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func CreateTable(db *sql.DB) error {

	createTable := "CREATE TABLE Metrics (type varchar(10), name varchar(50), value DOUBLE PRECISION, counter integer)"

	_, err := db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func AddGaugeMetric(db *sql.DB, nameOfMetric string, value string) error {
	logger.Log.Info("Add gauge metric into database begin")

	request := "MERGE INTO Metrics AS Metrics USING (SELECT type, name, value FROM Metrics WHERE type='gauge' AND name=$1) AS ExistMetric ON Metrics.type = ExistMetric.type AND Metrics.name = ExistMetric.name" +
		"WHEN NOT MATCHED THEN INSERT (type, name, value)  VALUES ('gauge', $2, $3)" +
		"WHEN MATCHED THEN UPDATE SET value = $4"
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}
	_, err = db.Exec(request, nameOfMetric, nameOfMetric, floatValue, floatValue)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}
	return nil
}

func AddCounterMetric(db *sql.DB, nameOfMetric string, value string) error {
	logger.Log.Info("Add counter metric into database begin")
	request := "MERGE INTO Metrics USING (SELECT type, name, counter FROM Metrics WHERE type='counter' AND name=$1) AS ExistMetric ON Metrics.type = ExistMetric.type AND Metrics.name = ExistMetric.name" +
		" WHEN NOT MATCHED THEN INSERT (type, name, counter)  VALUES ('gauge', $2, $3)" +
		" WHEN MATCHED THEN UPDATE SET counter = $4"
	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	int64Value++
	_, err = db.Exec(request, nameOfMetric, nameOfMetric, int64Value, int64Value)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	logger.Log.Info("Add counter metric into database end without errors")
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
