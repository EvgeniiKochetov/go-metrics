package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"net/http"

	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/metric"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/storage"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/models"
)

var memory storage.MemStorage

func init() {
	memory = storage.NewMemStorage()
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

func UpdateUseJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if (req.MType != "gauge") && (req.MType != "counter") {
		logger.Log.Debug("unsupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	switch req.MType {
	case "gauge":

		err := memory.ChangeGauge(req.ID, strconv.FormatFloat(*req.Value, 'f', -1, 64))
		if err != nil {
			logger.Log.Info("cannot change gauge", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "counter":

		err := memory.ChangeCounter(req.ID, strconv.FormatInt(*req.Delta, 10))
		if err != nil {
			logger.Log.Info("cannot change gauge", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		{
			http.Error(w, "Mistake in request! Wrong type metric", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	resp := models.Metrics{
		ID:    req.ID,
		MType: req.MType,
		Delta: req.Delta,
		Value: req.Value,
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Info("error encoding response", zap.Error(err))
		return
	}

}

func ValueUseJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if (req.MType != "gauge") && (req.MType != "counter") {
		logger.Log.Debug("unsupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	switch req.MType {
	case "gauge":

		res, ok := memory.GetMetricGauge(req.ID)
		if ok {
			resFloat, err := strconv.ParseFloat(res, 10)
			if err != nil {
				req.Value = &resFloat
			}
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	case "counter":

		res, ok := memory.GetMetricCounter(req.ID)
		if ok {
			resInt, err := strconv.ParseInt(res, 10, 10)
			if err != nil {
				req.Delta = &resInt
			}

		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	default:
		{
			http.Error(w, "Mistake in request! Wrong type metric", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	if err := enc.Encode(req); err != nil {
		logger.Log.Info("error encoding response", zap.Error(err))
		return
	}

}
