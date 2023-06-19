package main

import (
	"net/http"
	"testing"

	handlerclient "github.com/EvgeniiKochetov/go-metrics-tpl/internal/hadlerclient"
)

type want struct {
	code        int
	response    string
	contentType string
}

func Test_sendMetrics(t *testing.T) {
	client = &http.Client{}
	type args struct {
		typeOfMetric  string
		nameOfMetric  string
		valueOfMetric string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "", args: args{typeOfMetric: "counter", nameOfMetric: "someCounter", valueOfMetric: "10"}, wantErr: true},
		{name: "", args: args{typeOfMetric: "unknownTypeOfMetric", nameOfMetric: "test1", valueOfMetric: "10.1"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handlerclient.SendMetrics(client, serveraddr, tt.args.typeOfMetric, tt.args.nameOfMetric, tt.args.valueOfMetric); (err != nil) != tt.wantErr {
				t.Errorf("sendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}