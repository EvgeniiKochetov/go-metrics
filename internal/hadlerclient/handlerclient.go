package handlerclient

import (
	"net/http"
	"runtime"
)

func GetMetrics() *runtime.MemStats {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return m
}

func SendMetrics(client *http.Client, serveraddr, typeOfMetric, nameOfMetric, valueOfMetric string) error {

	reqURL := serveraddr + typeOfMetric + "/" + nameOfMetric + "/" + valueOfMetric
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	req.Header.Add("Content-Type", "text/plain")
	if err != nil {
		return err
	}
	//fmt.Println(req)
	response, err := client.Do(req)
	if err != nil {
		//fmt.Println(response)
		return err
	}
	defer response.Body.Close()
	return err
}
