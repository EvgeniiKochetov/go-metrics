package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/config"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/filestorage"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/gzip"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/handler"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx"
	"go.uber.org/zap"
)

var (
	flagRunAddr         string
	flagLogLevel        string
	flagStoreInterval   string
	flagFileStoragePath string
	flagRestore         string
	flagDatabase        string
)

func main() {

	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}

	connectToDatabase()
}

func parseFlags() {
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `user`, `psw`, `user`)
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.StringVar(&flagStoreInterval, "i", "2s", "store interval")
	flag.StringVar(&flagFileStoragePath, "f", "metrics-db.json", "storage path")
	flag.StringVar(&flagRestore, "r", "false", "restore")
	flag.StringVar(&flagDatabase, "d", ps, "configuration of SQL server")

	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		flagStoreInterval = envStoreInterval + "s"
	}
	if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
		//flagFileStoragePath = envStorePath
		flagFileStoragePath = "metrics-db.json"
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		flagRestore = envRestore
	}

	if envDatabase := os.Getenv("DATABASE_DSN"); envDatabase != "" {
		flagDatabase = envDatabase
	}
}

func run() error {

	r := chi.NewRouter()

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	r.Use(logger.RequestLogger, gzip.MyGzipHandle)
	r.NotFoundHandler()
	r.Route("/", func(r chi.Router) {

		r.Get("/", handler.AllMetrics)
		r.Get("/ping", handler.Ping)
		r.Post("/update/{typeMetric}/{metric}/{value}", handler.Update)
		r.Get("/value/counter/{metric}", handler.MetricCounter)
		r.Get("/value/gauge/{metric}", handler.MetricGauge)
		r.Post("/update/", handler.UpdateUseJSON)
		r.Post("/value/", handler.ValueUseJSON)
	})
	if flRestore, err := strconv.ParseBool(flagRestore); err == nil {
		fmt.Println("RESTORE " + flagRestore)
		if flRestore {
			handler.Memory.LoadStorage("metrics-db.json")
		}
	}
	go filestorage.SaveInFile(flagFileStoragePath, flagStoreInterval)

	return http.ListenAndServe(flagRunAddr, r)

}

func connectToDatabase() error {
	fmt.Println(flagDatabase)
	db, err := sql.Open("pgx", flagDatabase)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	config.GetInstance().SetDB(db)
	return err
}
