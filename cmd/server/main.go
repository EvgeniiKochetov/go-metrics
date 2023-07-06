package main

import (
	"flag"
	"fmt"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/tmp"
	"net/http"
	"os"
	"strconv"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/gzip"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/handler"
)

var (
	flagRunAddr         string
	flagLogLevel        string
	flagStoreInterval   string
	flagFileStoragePath string
	flagRestore         string
)

func main() {

	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.StringVar(&flagStoreInterval, "i", "300s", "store interval")
	flag.StringVar(&flagFileStoragePath, "f", "metrics-db.json", "storage path")
	flag.StringVar(&flagRestore, "r", "true", "restore")

	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		flagStoreInterval = envStoreInterval
	}
	if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
		flagFileStoragePath = envStorePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		flagRestore = envRestore
	}

}

func run() error {

	r := chi.NewRouter()

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	r.Use(logger.RequestLogger, gzip.MyGzipHandle)

	r.Route("/", func(r chi.Router) {

		r.Get("/", handler.AllMetrics)
		r.Post("/update/{typeMetric}/{metric}/{value}", handler.Update)
		r.Get("/value/counter/{metric}", handler.MetricCounter)
		r.Get("/value/gauge/{metric}", handler.MetricGauge)
		r.Post("/update/", handler.UpdateUseJSON)
		r.Post("/value/", handler.ValueUseJSON)
	})
	if flRestore, err := strconv.ParseBool(flagRestore); err != nil {
		fmt.Println("RESTORE " + flagRestore)
		if flRestore {
			handler.Memory.LoadStorage(flagFileStoragePath)
		}
	}
	go tmp.SaveInFile(flagFileStoragePath, flagStoreInterval)

	return http.ListenAndServe(flagRunAddr, r)

}
