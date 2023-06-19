package main

import (
	"flag"

	"net/http"

	"os"

	"github.com/go-chi/chi/v5"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/handler"
)

func main() {

	var flagRunAddr string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if err := run(flagRunAddr); err != nil {
		panic(err)
	}

}

func run(port string) error {

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {

		r.Get("/", handler.AllMetrics)
		r.Post("/update/{typeMetric}/{metric}/{value}", handler.Update)
		r.Get("/value/counter/{metric}", handler.MetricCounter)
		r.Get("/value/gauge/{metric}", handler.MetricGauge)

	})
	return http.ListenAndServe(port, r)

}
