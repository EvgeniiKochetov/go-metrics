package storage

import (
	"encoding/json"
	"fmt"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/models"
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
	slice := make([]models.Metrics, 0)

	for k, v := range m.metricsgauge {
		pointer := float64(v)
		slice = append(slice, models.Metrics{
			ID:    k,
			MType: "gauge",
			Delta: nil,
			Value: &pointer,
		})
	}

	for k, v := range m.metricscounter {
		pointer := int64(v)
		slice = append(slice, models.Metrics{
			ID:    k,
			MType: "gauge",
			Delta: &pointer,
			Value: nil,
		})
	}

	data, err := json.Marshal(slice)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0666)
}

func (m *MemStorage) LoadStorage(filename string) error {
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!RESTORE " + filename)
	data, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("Ошибка открытия файла")
		return err
	}

	slice := make([]models.Metrics, 0)
	json.Unmarshal(data, &slice)
	fmt.Println(slice)
	for _, v := range slice {
		if v.MType == "gauge" {
			m.ChangeGauge(v.ID, strconv.FormatFloat(*v.Value, 'f', -1, 64))
		} else {
			m.ChangeGauge(v.ID, strconv.FormatInt(*v.Delta, 10))
		}
	}
	return nil
}
