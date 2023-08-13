package handler

import (
	"encoding/json"
	"fmt"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/config"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/database"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/models"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/storage"
)

var Memory storage.MemStorage

func init() {
	Memory = storage.NewMemStorage()
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
		err = Memory.ChangeGauge(nameOfMetric, valueOfMetric)
	case "counter":
		err = Memory.ChangeCounter(nameOfMetric, valueOfMetric)
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

	slice, ok := Memory.GetAllMetrics()
	if !ok {
		w.Write([]byte("No data"))
	}
	for _, v := range slice {
		w.Write([]byte(v + "\n"))
	}

}

func MetricGauge(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	res, ok := Memory.GetMetricGauge(metric)
	if ok {
		w.Write([]byte(res))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("No value"))

}

func MetricCounter(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	res, ok := Memory.GetMetricCounter(metric)
	if ok {
		w.Write([]byte(res))
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)

}

func UpdateUseJSON(w http.ResponseWriter, r *http.Request) {

	fmt.Println()
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
	var value string
	switch req.MType {
	case "gauge":

		fmt.Println(*req.Value)

		if float64(*req.Value) == float64(int(*req.Value)) {
			value = strconv.FormatFloat(float64(*req.Value), 'f', 1, 64)
		} else {
			value = strconv.FormatFloat(*req.Value, 'f', 24, 64)
		}
		fmt.Println(value)
		err := Memory.ChangeGauge(req.ID, value)
		if err != nil {
			logger.Log.Info("cannot change gauge", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		db := config.GetInstance().GetDatabaseConnection()
		if db != nil {
			database.AddGaugeMetric(db, req.ID, value)
		}

	case "counter":

		err := Memory.ChangeCounter(req.ID, strconv.FormatInt(*req.Delta, 10))
		if err != nil {
			logger.Log.Info("cannot change gauge", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		db := config.GetInstance().GetDatabaseConnection()
		if db != nil {
			database.AddCounterMetric(db, req.ID, strconv.FormatInt(*req.Delta, 10))
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
	w.WriteHeader(http.StatusOK)
}

func ValueUseJSON(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logger.Log.Info("got request with bad method", zap.String("method", r.Method))
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
		logger.Log.Info("unsupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	switch req.MType {
	case "gauge":

		res, ok := Memory.GetMetricGauge(req.ID)

		if ok {
			resFloat, err := strconv.ParseFloat(res, 64)

			if err != nil {
				logger.Log.Info("unsupported request type")
				w.WriteHeader(http.StatusNotFound)

			} else {
				req.Value = &resFloat
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case "counter":

		res, ok := Memory.GetMetricCounter(req.ID)

		if ok {
			resInt, err := strconv.ParseInt(res, 10, 64)
			if err != nil {
				fmt.Println("Ошибка конвертации counter", req.ID, resInt)
				w.WriteHeader(http.StatusNotFound)
			} else {
				req.Delta = &resInt
			}

		} else {
			fmt.Println("Ошибка поиска  counter", req.ID)
			w.WriteHeader(http.StatusNotFound)
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
	w.WriteHeader(http.StatusOK)
}

func UpdatesUseJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req, resp []models.Metrics

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, v := range req {
		if (v.MType != "gauge") && (v.MType != "counter") {
			logger.Log.Debug("unsupported request type", zap.String("type", v.MType))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		var value string
		switch v.MType {
		case "gauge":

			fmt.Println(*v.Value)

			if float64(*v.Value) == float64(int(*v.Value)) {
				value = strconv.FormatFloat(float64(*v.Value), 'f', 1, 64)
			} else {
				value = strconv.FormatFloat(*v.Value, 'f', 24, 64)
			}
			fmt.Println(value)
			err := Memory.ChangeGauge(v.ID, value)
			if err != nil {
				logger.Log.Info("cannot change gauge", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			db := config.GetInstance().GetDatabaseConnection()
			if db != nil {
				database.AddGaugeMetric(db, v.ID, value)
			}

		case "counter":

			err := Memory.ChangeCounter(v.ID, strconv.FormatInt(*v.Delta, 10))
			if err != nil {
				logger.Log.Info("cannot change gauge", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			db := config.GetInstance().GetDatabaseConnection()
			if db != nil {
				database.AddCounterMetric(db, v.ID, strconv.FormatInt(*v.Delta, 10))
			}

		default:
			{
				http.Error(w, "Mistake in request! Wrong type metric", http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Info("error encoding response", zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)

}

func Ping(w http.ResponseWriter, r *http.Request) {
	db := config.GetInstance().GetDatabaseConnection()

	fmt.Println("Ping start ")
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		fmt.Println("Ping end: ", w.Header())
		return
	}
	w.WriteHeader(http.StatusOK)
}
