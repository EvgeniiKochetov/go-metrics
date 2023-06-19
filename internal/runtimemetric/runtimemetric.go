package runtimemetric

import (
	"fmt"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/storage"

	"net/http"

	"reflect"

	"time"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/hadlerclient"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/metric"
)

func Run(client *http.Client, serveraddr string, reportInterval, pollInterval int) {
	//time.Sleep(time.Second * time.Duration(reportInterval))
	metricmap := metric.GetMapMetrics()
	storageMetric := storage.NewMemStorage()
	var value string
	var ok bool

	for {
		m := handlerclient.GetMetrics()
		val := reflect.ValueOf(m).Elem()
		typeOfT := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			nameMetric := typeOfT.Field(i).Name
			if nameMetric == "counter" {
				storageMetric.ChangeGauge(nameMetric, fmt.Sprint(field.Interface()))
			} else {
				storageMetric.ChangeCounter(nameMetric, fmt.Sprint(field.Interface()))
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
			}
			time.Sleep(time.Second * time.Duration(reportInterval))
		}

		time.Sleep(time.Second * time.Duration(pollInterval))
	}

}
