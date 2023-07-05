package main

import (
	"flag"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/gzip"
	"go.uber.org/zap"

	"net/http"

	"os"

	"github.com/go-chi/chi/v5"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/handler"
)

var (
	flagRunAddr  string
	flagLogLevel string
)

func main() {

	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")

	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {

	r := chi.NewRouter()

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	r.Use(logger.RequestLogger, gzip.MyGzipHandle)
	//r.Use(logger.RequestLogger, gzip.MyGzipMiddleware)
	//r.Use(logger.RequestLogger)

	r.Route("/", func(r chi.Router) {

		r.Get("/", handler.AllMetrics)
		r.Post("/update/{typeMetric}/{metric}/{value}", handler.Update)
		r.Get("/value/counter/{metric}", handler.MetricCounter)
		r.Get("/value/gauge/{metric}", handler.MetricGauge)
		r.Post("/update/", handler.UpdateUseJSON)
		r.Post("/value/", handler.ValueUseJSON)
	})
	return http.ListenAndServe(flagRunAddr, r)

}
