package runtimemetric

import (
	"fmt"
	"math/rand"
	"strconv"

	"net/http"

	"reflect"

	"time"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/hadlerclient"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/metric"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/storage"
)

func Run(client *http.Client, serveraddr string, reportInterval, pollInterval int) {

	metricmap := metric.GetMapMetrics()
	storageMetric := storage.NewMemStorage()
	var value string
	var ok bool
	var counter int

	for {
		m := handlerclient.GetMetrics()
		val := reflect.ValueOf(m).Elem()
		typeOfT := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			nameMetric := typeOfT.Field(i).Name
			if nameMetric == "counter" {
				storageMetric.ChangeCounter(nameMetric, fmt.Sprint(field.Interface()))
			} else {
				storageMetric.ChangeGauge(nameMetric, fmt.Sprint(field.Interface()))
			}
		}

		for k, typeOfMetric := range metricmap {
			if typeOfMetric == "counter" {
				value, ok = storageMetric.GetMetricCounter(k)
			} else {
				value, ok = storageMetric.GetMetricGauge(k)
			}

			if ok {
				handlerclient.SendMetrics(client, serveraddr, typeOfMetric, k, value)
				counter++

				handlerclient.SendMetrics(client, serveraddr, "counter", "PollCount", strconv.FormatInt(int64(counter), 10))
				handlerclient.SendMetrics(client, serveraddr, "gauge", "RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64))
			}
			time.Sleep(time.Second * time.Duration(reportInterval))
		}

		time.Sleep(time.Second * time.Duration(pollInterval))
	}

}
