package internal

import "time"

type CostBreakdown struct {
	CPUCost     float64 `json:"cpu_cost"`
	MemoryCost  float64 `json:"memory_cost"`
	StorageCost float64 `json:"storage_cost"`
	NetworkCost float64 `json:"network_cost"`
	TotalCost   float64 `json:"totalCost"`
}

type CostOverview struct {
	TotalCost      CostBreakdown   `json:"total_cost"`
	NamespacesCost []NamespaceCost `json:"namespace_costs"`
	Timestamp      time.Time       `json:"timestamp"`
}

type NamespaceCost struct {
	Namespace   string    `json:"namespace"`
	CPUCost     float64   `json:"cpu_cost"`
	MemoryCost  float64   `json:"memory_cost"`
	StorageCost float64   `json:"storage_cost"`
	NetworkCost float64   `json:"network_cost"`
	TotalCost   float64   `json:"total_cost"`
	PodCount    int       `json:"pod_count"`
	Pods        []PodCost `json:"pods,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

type PodCost struct {
	Name        string    `json:"name"`
	Namespace   string    `json:"namespace"`
	CPUCost     float64   `json:"cpu_cost"`
	MemoryCost  float64   `json:"memory_cost"`
	StorageCost float64   `json:"storage_cost"`
	NetworkCost float64   `json:"network_cost"`
	TotalCost   float64   `json:"total_cost"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage int64     `json:"memory_usage"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Timestamp   time.Time `json:"timestamp"`
}

type NodeCost struct {
	Name           string    `json:"name"`
	CPUCost        float64   `json:"cpu_cost"`
	MemoryCost     float64   `json:"memory_cost"`
	StorageCost    float64   `json:"storage_cost"`
	TotalCost      float64   `json:"total_cost"`
	CPUCapacity    string    `json:"cpu_capacity"`
	MemoryCapacity string    `json:"memory_capacity"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    int64     `json:"memory_usage"`
	PodCount       int       `json:"pod_count"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
}

type CostHistory struct {
	Period    string             `json:"period"`
	StartTime time.Time          `json:"start_time"`
	EndTime   time.Time          `json:"end_time"`
	Data      []CostHistoryPoint `json:"data"`
}

type CostHistoryPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUCost     float64   `json:"cpu_cost"`
	MemoryCost  float64   `json:"memory_cost"`
	StorageCost float64   `json:"storage_cost"`
	NetworkCost float64   `json:"network_cost"`
	TotalCost   float64   `json:"total_cost"`
}

type ResourceUsage struct {
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    int64     `json:"memory_usage"`
	StorageUsage   int64     `json:"storage_usage"`
	NetworkRxBytes float64   `json:"network_rx_bytes"`
	NetworkTxBytes float64   `json:"network_tx_bytes"`
	Timestamp      time.Time `json:"timestamp"`
}

type ClusterMetrics struct {
	TotalNodes         int           `json:"total_nodes"`
	TotalPods          int           `json:"total_pods"`
	TotalNamespaces    int           `json:"total_namespaces"`
	CPUUtilization     float64       `json:"cpu_utilization"`
	MemoryUtilization  float64       `json:"memory_utilization"`
	StorageUtilization float64       `json:"storage_utilization"`
	ResourceUsage      ResourceUsage `json:"resource_usage"`
	Timestamp          time.Time     `json:"timestamp"`
}

type PromMetrics struct {
	CPUMetrics     []MetricPoint `json:"cpu_metrics"`
	MemoryMetrics  []MetricPoint `json:"memory_metrics"`
	NetworkMetrics []MetricPoint `json:"network_metrics"`
	StorageMetrics []MetricPoint `json:"storage_metrics"`
	Timestamp      time.Time     `json:"timestamp"`
}

type MetricPoint struct {
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels"`
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}
