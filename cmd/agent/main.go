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

var (
	flagServPort       string
	flagReportInterval string
	flagPollInterval   string
	key                string
	flagRateLimit      string
)

func main() {

	flag.StringVar(&flagServPort, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagReportInterval, "r", "10", "report interval")
	flag.StringVar(&flagPollInterval, "p", "2", "poll interval")
	flag.StringVar(&key, "k", "", "key for hash")
	flag.StringVar(&flagRateLimit, "l", "", "rate limit")
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

	if envKey := os.Getenv("KEY"); envKey != "" {
		key = envKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		flagRateLimit = envRateLimit
	}

	serveraddr = "http://" + flagServPort + "/update/"
	pollInterval, err := strconv.Atoi(flagPollInterval)
	if err != nil {
		panic(err)
	}

	runtimemetric.Run(client, serveraddr, reportInterval, pollInterval, key)

}
