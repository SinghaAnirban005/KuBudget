package metrics

import (
	"net/http"
	"time"
)

type PrometheusClient struct {
	baseUrl string
	Client  *http.Client
}

func NewPrometheusClient(baseUrl string) *PrometheusClient {
	return &PrometheusClient{
		baseUrl: baseUrl,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
