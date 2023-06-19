package runtimemetric

import (
	"fmt"

	"net/http"

	"reflect"

	"time"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/hadlerclient"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/metric"
)

func Run(client *http.Client, serveraddr string, reportInterval, pollInterval int) {
	//time.Sleep(time.Second * time.Duration(reportInterval-pollInterval))
	var metricmap map[string]string
	metricmap = metric.GetMapMetrics()

	for {
		m := handlerclient.GetMetrics()
		val := reflect.ValueOf(m).Elem()
		typeOfT := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			typeOfMetric, ok := metricmap[typeOfT.Field(i).Name]
			if ok {

				handlerclient.SendMetrics(client, serveraddr, typeOfMetric, typeOfT.Field(i).Name, fmt.Sprint(field.Interface()))
			}

		}
		time.Sleep(time.Second * time.Duration(pollInterval))
	}

}
