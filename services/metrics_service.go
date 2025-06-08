package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SinghaAnirban005/KuBudget/internal"
	"github.com/SinghaAnirban005/KuBudget/pkg/kubernetes"
	"github.com/SinghaAnirban005/KuBudget/pkg/prometheus"
)

type MetricsService struct {
	k8sClient  *kubernetes.Client
	promClient *prometheus.Client
}

func NewMetricsService(k8sClient *kubernetes.Client, promClient *prometheus.Client) *MetricsService {
	return &MetricsService{
		k8sClient:  k8sClient,
		promClient: promClient,
	}
}

func (s *MetricsService) GetPrometheusMetrics(ctx context.Context, namespace, pod string) (*internal.PromMetrics, error) {
	now := time.Now()

	cpuQuery := "rate(container_cpu_usage_seconds_total[5m])"
	if namespace != "" {
		cpuQuery = fmt.Sprintf(`rate(container_cpu_usage_seconds_total{namespace="%s"}[5m])`, namespace)
	}
	if pod != "" {
		cpuQuery = fmt.Sprintf(`rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s"}[5m])`, namespace, pod)
	}

	cpuResult, err := s.promClient.Query(ctx, cpuQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU metrics: %v", err)
	}

	memQuery := "container_memory_usage_bytes"
	if namespace != "" {
		memQuery = fmt.Sprintf(`container_memory_usage_bytes{namespace="%s"}`, namespace)
	}
	if pod != "" {
		memQuery = fmt.Sprintf(`container_memory_usage_bytes{namespace="%s",pod="%s"}`, namespace, pod)
	}

	memResult, err := s.promClient.Query(ctx, memQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory metrics: %v", err)
	}

	netRxQuery := "rate(container_network_receive_bytes_total[5m])"
	netTxQuery := "rate(container_network_transmit_bytes_total[5m])"
	if namespace != "" {
		netRxQuery = fmt.Sprintf(`rate(container_network_receive_bytes_total{namespace="%s"}[5m])`, namespace)
		netTxQuery = fmt.Sprintf(`rate(container_network_transmit_bytes_total{namespace="%s"}[5m])`, namespace)
	}

	netRxResult, _ := s.promClient.Query(ctx, netRxQuery)
	netTxResult, _ := s.promClient.Query(ctx, netTxQuery)

	cpuMetrics := s.convertToMetricPoints(cpuResult)
	memMetrics := s.convertToMetricPoints(memResult)
	networkMetrics := append(s.convertToMetricPoints(netRxResult), s.convertToMetricPoints(netTxResult)...)

	return &internal.PromMetrics{
		CPUMetrics:     cpuMetrics,
		MemoryMetrics:  memMetrics,
		NetworkMetrics: networkMetrics,
		StorageMetrics: []internal.MetricPoint{},
		Timestamp:      now,
	}, nil
}

func (s *MetricsService) GetClusterMetrics(ctx context.Context) (*internal.ClusterMetrics, error) {
	nodes, err := s.k8sClient.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("Failed to get nodes %v", err)
	}

	namespaces, err := s.k8sClient.GetNamespaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %v", err)
	}

	totalPods := 0
	for _, ns := range namespaces.Items {
		pods, err := s.k8sClient.GetPods(ns.Name)
		if err != nil {
			continue
		}
		totalPods += len(pods.Items)
	}

	cpuQuery := `(1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m]))) * 100`
	cpuResult, err := s.promClient.Query(ctx, cpuQuery)
	fmt.Println("Hi there")
	var cpuUtilization float64
	if err == nil && len(cpuResult.Data.Result) > 0 {
		cpuUtilization = cpuResult.Data.Result[0].Value
	}

	memQuery := `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`
	memResult, err := s.promClient.Query(ctx, memQuery)
	var memUtilization float64
	if err == nil && len(memResult.Data.Result) > 0 {
		memUtilization = float64(memResult.Data.Result[0].Value)
	}

	resourceUsage, err := s.GetResourceUsage(ctx, "", "")
	if err != nil {
		resourceUsage = &internal.ResourceUsage{
			Timestamp: time.Now(),
		}
	}

	return &internal.ClusterMetrics{
		TotalNodes:         len(nodes.Items),
		TotalPods:          totalPods,
		TotalNamespaces:    len(namespaces.Items),
		CPUUtilization:     cpuUtilization,
		MemoryUtilization:  memUtilization,
		StorageUtilization: 0,
		ResourceUsage:      *resourceUsage,
		Timestamp:          time.Now(),
	}, nil

}

func (s *MetricsService) GetResourceUsage(ctx context.Context, namespace, pod string) (*internal.ResourceUsage, error) {
	cpuUsage, err := s.promClient.GetCPUUsage(ctx, namespace, pod)
	if err != nil {
		cpuUsage = 0
	}

	memoryUsage, err := s.promClient.GetMemoryUsage(ctx, namespace, pod)
	if err != nil {
		memoryUsage = 0
	}

	rxBytes, txBytes, err := s.promClient.GetNetworkIO(ctx, namespace, pod)
	if err != nil {
		rxBytes = 0
		txBytes = 0
	}

	return &internal.ResourceUsage{
		CPUUsage:       cpuUsage,
		MemoryUsage:    int64(memoryUsage),
		StorageUsage:   0,
		NetworkRxBytes: rxBytes,
		NetworkTxBytes: txBytes,
		Timestamp:      time.Now(),
	}, nil
}

func (s *MetricsService) convertToMetricPoints(result *prometheus.QueryResult) []internal.MetricPoint {
	if result == nil || len(result.Data.Result) == 0 {
		return []internal.MetricPoint{}
	}

	points := make([]internal.MetricPoint, len(result.Data.Result))
	for i, sample := range result.Data.Result {
		labels := make(map[string]string)

		points[i] = internal.MetricPoint{
			Timestamp: sample.Timestamp,
			Value:     float64(sample.Value),
			Labels:    labels,
		}

	}
	return points
}
