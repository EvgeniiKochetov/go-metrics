package runtimemetric

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/mem"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/hadlerclient"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/metric"
)

type metricData struct {
	typeOfMetric string
	nameOfMetric string
	value        string
}

func GetExtraMetrics(pollInterval int, metrics chan metricData) {
	for {
		v, _ := mem.VirtualMemory()

		metrics <- metricData{"gauge", "TotalMemory", string(v.Total)}
		metrics <- metricData{"gauge", "FreeMemory", string(v.Free)}
		metrics <- metricData{"gauge", "Used", string(v.Used)}
		time.Sleep(time.Second * time.Duration(pollInterval))
	}
}

func GetMetrics(pollInterval int, metrics chan metricData) {

	for {
		m := handlerclient.GetMetrics()
		val := reflect.ValueOf(m).Elem()
		typeOfT := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			nameMetric := typeOfT.Field(i).Name
			metrics <- metricData{"gauge", nameMetric, fmt.Sprint(field.Interface())}
		}

		time.Sleep(time.Second * time.Duration(pollInterval))
	}
}

func SendMetrics(reportInterval int, client *http.Client, serveraddr string, key string, metrics chan metricData) {
	for {

		var metricdata metricData
		metricdata = <-metrics
		var counter = 0
		err := handlerclient.SendMetric(client, serveraddr, metricdata.typeOfMetric, metricdata.nameOfMetric, metricdata.value, key)
		if err != nil {
			fmt.Println("error when send metric: ", err)
			return
		}
		err = handlerclient.SendMetric(client, serveraddr, "counter", "PollCount", strconv.FormatInt(int64(counter), 10), key)
		if err != nil {
			fmt.Println("error when send metric count: ", err)
			return
		}
		//fmt.Println("get metric: ", metricdata)
		time.Sleep(time.Second * time.Duration(reportInterval))
	}
}

func Run(client *http.Client, serveraddr string, reportInterval, pollInterval int, key string, rateLimit string) {

	metrics := make(chan metricData, len(metric.GetMapMetrics()))

	go GetExtraMetrics(pollInterval, metrics)
	go GetMetrics(pollInterval, metrics)

	defer close(metrics)
	limit, _ := strconv.Atoi(rateLimit)
	for a := 1; a <= limit; a++ {
		go SendMetrics(reportInterval, client, serveraddr, key, metrics)
	}
	for {

	}
	//metricmap := metric.GetMapMetrics()
	//storageMetric := storage.NewMemStorage()
	//var value string
	//var ok bool
	//var counter int

	//for {
	//	m := handlerclient.GetMetrics()
	//	val := reflect.ValueOf(m).Elem()
	//	typeOfT := val.Type()
	//
	//	for i := 0; i < val.NumField(); i++ {
	//		field := val.Field(i)
	//		nameMetric := typeOfT.Field(i).Name
	//		if nameMetric == "counter" {
	//			storageMetric.ChangeCounter(nameMetric, fmt.Sprint(field.Interface()))
	//		} else {
	//			storageMetric.ChangeGauge(nameMetric, fmt.Sprint(field.Interface()))
	//		}
	//	}
	//
	//	for k, typeOfMetric := range metricmap {
	//		if typeOfMetric == "counter" {
	//			value, ok = storageMetric.GetMetricCounter(k)
	//		} else {
	//			value, ok = storageMetric.GetMetricGauge(k)
	//		}
	//
	//		if ok {
	//			handlerclient.SendMetric(client, serveraddr, typeOfMetric, k, value, key)
	//			counter++
	//
	//			handlerclient.SendMetric(client, serveraddr, "counter", "PollCount", strconv.FormatInt(int64(counter), 10), key)
	//			handlerclient.SendMetric(client, serveraddr, "gauge", "RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64), key)
	//		}
	//		time.Sleep(time.Second * time.Duration(reportInterval))
	//	}
	//
	//	time.Sleep(time.Second * time.Duration(pollInterval))
	//}

}
