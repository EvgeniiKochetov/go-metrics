package main

import (
	"fmt"
	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type want struct {
	code        int
	response    string
	contentType string
}

func TestUpdate(t *testing.T) {

	testCases := []struct {
		handlerName  string
		method       string
		request      string
		expectedCode int
	}{
		{handlerName: "update", method: http.MethodPost, request: "http://localhost:8080/update/gauge/param/param/12.0", expectedCode: http.StatusBadRequest},
		{handlerName: "allMetrics", method: http.MethodGet, request: "http://localhost:8080/", expectedCode: http.StatusOK},
		{handlerName: "metricGauge", method: http.MethodGet, request: "http://localhost:8080/value/uknownmetric/", expectedCode: http.StatusNotFound},
	}
	var r *http.Request
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {

			w := httptest.NewRecorder()
			r = httptest.NewRequest(tc.method, tc.request, nil)
			fmt.Println(tc.request)
			switch tc.handlerName {
			case "update":
				handler.Update(w, r)
			case "allMetrics":
				handler.AllMetrics(w, r)
			case "metricGauge":
				handler.MetricGauge(w, r)
			case "metricCounter":
				handler.MetricCounter(w, r)
			}

			if tc.expectedCode != 0 {
				assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			}
		})
	}
}
