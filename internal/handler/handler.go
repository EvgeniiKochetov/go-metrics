package handler

import (
	"encoding/json"

	"net/http"

	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/metric"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/storage"
)

var memory storage.MemStorage

func Init() {
	memory = storage.MemStorage{}
}

func Update(w http.ResponseWriter, r *http.Request) {
	var err error

	typeOfMetric := chi.URLParam(r, "typeMetric")
	nameOfMetric := chi.URLParam(r, "metric")
	valueOfMetric := chi.URLParam(r, "value")

	underparts := strings.Split(r.URL.Path, "/")
	if len(underparts) < 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch typeOfMetric {
	case "gauge":
		err = memory.ChangeGauge(nameOfMetric, valueOfMetric)
	case "counter":
		err = memory.ChangeCounter(nameOfMetric, valueOfMetric)
	default:
		{
			http.Error(w, "Mistake in request! Wrong type metric", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if err != nil {
		http.Error(w, "Mistake in request! Wrong numbers", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	}

	response, _ := json.Marshal(`{"status":"ok"}`)
	w.Write(response)
}

func AllMetrics(w http.ResponseWriter, _ *http.Request) {

	slice, ok := metric.StorageMetric.GetAllMetrics()
	if !ok {
		w.Write([]byte("No data"))
	}
	for _, v := range slice {
		w.Write([]byte(v + "\n"))
	}

}

func MetricGauge(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	res, ok := memory.GetMetricGauge(metric)
	if ok {
		w.Write([]byte(res))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("No value"))

}

func MetricCounter(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	res, ok := memory.GetMetricCounter(metric)
	if ok {
		w.Write([]byte(res))
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)

}
