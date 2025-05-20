package models

import "time"

type NodeMetrics struct {
	Name           string    `json:"name"`
	CPUUsage       float64   `json:"cpuUsage"`
	CPUCapacity    float64   `json:"cpuCapacity"`
	MemoryUsage    float64   `json:"memoryUsage"`
	MemoryCapacity float64   `json:"memoryCapacity"`
	PodCount       int       `json:"podCount"`
	Cost           float64   `json:"cost"`
	Timestamp      time.Time `json:"timestamp"`
}

type PodMetrics struct {
	Name        string    `json:"name"`
	Namespace   string    `json:"namespace"`
	CPUUsage    float64   `json:"cpuUsage"`
	CPULimit    float64   `json:"cpuLimit"`
	MemoryUsage float64   `json:"memoryUsage"`
	MemoryLimit float64   `json:"memoryLimit"`
	Cost        float64   `json:"cost"`
	Timestamp   time.Time `json:"timestamp"`
	NodeName    string    `json:"nodeName"`
}

type ResourceList struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type PrometheusQueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string                  `json:"resultType"`
		Result     []PrometheusQueryResult `json:"result"`
	} `json:"data"`
}

type MetricsServerPod struct {
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Timestamp  time.Time `json:"timestamp"`
	Window     string    `json:"window"`
	Containers []struct {
		Name  string       `json:"name"`
		Usage ResourceList `json:"usage"`
	} `json:"containers"`
}

type MetricsServerResponse struct {
	Items []MetricsServerPod `json:"items"`
}
