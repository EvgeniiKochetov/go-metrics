package metric

import "github.com/EvgeniiKochetov/go-metrics-tpl/internal/storage"

var StorageMetric storage.MemStorage

func GetMapMetrics() map[string]string {
	return map[string]string{
		"Alloc":         "gauge",
		"BuckHashSys":   "gauge",
		"Frees":         "gauge",
		"GCCPUFraction": "gauge",
		"GCSys":         "gauge",
		"HeapAlloc":     "gauge",
		"HeapIdle":      "gauge",
		"HeapInuse":     "gauge",
		"HeapObjects":   "gauge",
		"HeapReleased":  "gauge",
		"HeapSys":       "gauge",
		"LastGC":        "gauge",
		"Lookups":       "gauge",
		"MCacheInuse":   "gauge",
		"MCacheSys":     "gauge",
		"MSpanInuse":    "gauge",
		"MSpanSys":      "gauge",
		"Mallocs":       "gauge",
		"NextGC":        "gauge",
		"NumForcedGC":   "gauge",
		"NumGC":         "gauge",
		"OtherSys":      "gauge",
		"PauseTotalNs":  "gauge",
		"StackInuse":    "gauge",
		"StackSys":      "gauge",
		"Sys":           "gauge",
		"TotalAlloc":    "gauge",
		"PollCount":     "counter",
		"RandomValue ":  "gauge",
	}
}
