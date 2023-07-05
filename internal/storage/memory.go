package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type gauge float64
type counter int64

type MemStorage struct {
	metricsgauge   map[string]gauge
	metricscounter map[string]counter
}

func (m *MemStorage) ChangeGauge(name string, value string) error {
	res, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.metricsgauge[name] = gauge(res)
	return nil
}

func (m *MemStorage) ChangeCounter(name string, value string) error {
	res, err := strconv.Atoi(value)
	if err == nil {
		m.metricscounter[name] += counter(res)
	}
	return err
}

func NewMemStorage() MemStorage {
	p := MemStorage{metricsgauge: map[string]gauge{}, metricscounter: map[string]counter{}}
	return p
}
func (m *MemStorage) GetMetricGauge(name string) (string, bool) {
	k, ok := m.metricsgauge[name]

	return strconv.FormatFloat(float64(k), 'f', -1, 64), ok

}

func (m *MemStorage) GetMetricCounter(name string) (string, bool) {
	k, ok := m.metricscounter[name]

	return strconv.FormatInt(int64(k), 10), ok

}

func (m *MemStorage) GetAllMetrics() ([]string, bool) {
	if len(m.metricsgauge) == 0 && len(m.metricscounter) == 0 {
		return nil, false
	}
	var slice []string
	for _, v := range m.metricsgauge {
		slice = append(slice, "gauge: ", strconv.FormatFloat(float64(v), 'f', -1, 64))
	}

	for _, v := range m.metricscounter {
		slice = append(slice, "counter: ", strconv.FormatFloat(float64(v), 'f', -1, 64))
	}

	return slice, true
}

func (m *MemStorage) SaveStorage(filename string) error {
	slice, _ := m.GetAllMetrics()
	data, err := json.Marshal(slice)
	fmt.Println(string(data))
	if err != nil {
		return err
	}
	fmt.Println(data)
	return os.WriteFile(filename, data, 0666)
}
