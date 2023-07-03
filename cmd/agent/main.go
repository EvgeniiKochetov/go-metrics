package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/metric"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/runtimemetric"
)

var reportInterval int
var pollInterval int
var client *http.Client
var serveraddr string
var metrics map[string]string

func init() {
	metrics = metric.GetMapMetrics()
	client = &http.Client{}
}

func main() {
	var flagReportInterval string
	var flagPollInterval string
	var flagServPort string

	flag.StringVar(&flagServPort, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagReportInterval, "r", "10", "report interval")
	flag.StringVar(&flagPollInterval, "p", "2", "poll interval")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagServPort = envRunAddr
	}

	if envRunAddr := os.Getenv("REPORT_INTERVAL"); envRunAddr != "" {
		flagReportInterval = envRunAddr
	}

	if envRunAddr := os.Getenv("POLL_INTERVAL"); envRunAddr != "" {
		flagPollInterval = envRunAddr
	}

	serveraddr = "http://" + flagServPort + "/update/"
	pollInterval, err := strconv.Atoi(flagPollInterval)
	if err != nil {
		panic(err)
	}

	runtimemetric.Run(client, serveraddr, reportInterval, pollInterval)

}
