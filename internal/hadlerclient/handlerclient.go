package handlerclient

import (
	"crypto/sha256"
	"net/http"
	"runtime"
)

func GetMetrics() *runtime.MemStats {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return m
}

func SendMetric(client *http.Client, serveraddr, typeOfMetric, nameOfMetric, valueOfMetric string, key string) error {

	reqURL := serveraddr + typeOfMetric + "/" + nameOfMetric + "/" + valueOfMetric
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	req.Header.Add("Content-Type", "text/plain")
	if key != "" {
		req.Header.Add("HashSHA256", key)
	}

	if err != nil {
		return err
	}

	response, err := client.Do(req)
	if err != nil {

		return err
	}
	defer response.Body.Close()
	return err
}

func GetHash(src []byte) []byte {
	// создаём новый hash.Hash, вычисляющий контрольную сумму SHA-256
	h := sha256.New()
	// передаём байты для хеширования
	h.Write(src)
	// вычисляем хеш
	dst := h.Sum(nil)
	return dst
}
