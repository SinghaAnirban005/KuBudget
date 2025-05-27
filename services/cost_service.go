package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/SinghaAnirban005/KuBudget/internal"
	"github.com/SinghaAnirban005/KuBudget/pkg/kubernetes"
	"github.com/SinghaAnirban005/KuBudget/pkg/prometheus"
)

type CostService struct {
	k8sClient  *kubernetes.Client
	promClient *prometheus.Client
	config     *internal.Config
}

func NewCostService(k8sClient *kubernetes.Client, promClient *prometheus.Client) *CostService {
	return &CostService{
		k8sClient:  k8sClient,
		promClient: promClient,
		config:     internal.LoadConfig(),
	}
}

func (s *CostService) GetCostOverview(ctx context.Context) (*internal.CostOverview, error) {
	namespaces, err := s.k8sClient.GetNamespaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %v", err)
	}

	var totalCPUCost, totalMemoryCost, totalStorageCost float64
	namespaceCosts := make([]internal.NamespaceCost, 0, len(namespaces.Items))

	for _, ns := range namespaces.Items {
		nsCost, err := s.calculateNamespaceCost(ctx, ns.Name)
		if err != nil {
			continue
		}

		totalCPUCost += nsCost.CPUCost
		totalMemoryCost += nsCost.MemoryCost
		totalStorageCost += nsCost.StorageCost

		namespaceCosts = append(namespaceCosts, *nsCost)
	}

	return &internal.CostOverview{
		TotalCost: internal.CostBreakdown{
			CPUCost:     totalCPUCost,
			MemoryCost:  totalMemoryCost,
			StorageCost: totalStorageCost,
		},
		NamespacesCost: namespaceCosts,
		Timestamp:      time.Now(),
	}, nil

}

func (s *CostService) GetCostHistory(ctx context.Context, duration time.Duration, step time.Duration, namespace string) (*internal.CostHistory, error) {
	endTime := time.Now()
	startTime := endTime.Add(-duration)

	cpuQuery := "rate(container_cpu_usage_seconds_total[5m])"
	if namespace != "" {
		cpuQuery = fmt.Sprintf(`rate(container_cpu_usage_seconds_total{namespace="%s"}[5m])`, namespace)
	}

	cpuResult, err := s.promClient.QueryRange(ctx, cpuQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU history: %v", err)
	}

	memQuery := "container_memory_usage_bytes"
	if namespace != "" {
		memQuery = fmt.Sprintf(`container_memory_usage_bytes{namespace="%s"}`, namespace)
	}

	memResult, err := s.promClient.QueryRange(ctx, memQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory history: %v", err)
	}

	historyPoints := s.convertToCostHistory(*cpuResult, *memResult)

	return &internal.CostHistory{
		Period:    duration.String(),
		StartTime: startTime,
		EndTime:   endTime,
		Data:      historyPoints,
	}, nil
}

func (s *CostService) GetNodeCosts(ctx context.Context) ([]internal.NodeCost, error) {
	nodes, err := s.k8sClient.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %v", err)
	}

	costs := make([]internal.NodeCost, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		cost, err := s.calculateNodeCost(ctx, node.Namespace, "", node.Name)
		if err != nil {
			continue
		}
		costs = append(costs, *cost)
	}

	return costs, err
}

func (s *CostService) GetNamespaceCosts(ctx context.Context) ([]internal.NamespaceCost, error) {
	namespaces, err := s.k8sClient.GetNamespaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %v", err)
	}

	costs := make([]internal.NamespaceCost, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		cost, err := s.calculateNamespaceCost(ctx, ns.Name)
		if err != nil {
			continue
		}
		costs = append(costs, *cost)
	}

	return costs, nil

}

func (s *CostService) calculateNodeCost(ctx context.Context, namespace, pod, nodeName string) (*internal.NodeCost, error) {
	cpuUsage, err := s.promClient.GetCPUUsage(ctx, namespace, pod)
	if err != nil {
		cpuUsage = 0
	}

	cpuCost := cpuUsage * s.config.CPUCostPerHour

	memoryUsage, err := s.promClient.GetMemoryUsage(ctx, namespace, pod)
	if err != nil {
		memoryUsage = 0
	}
	memoryGB := memoryUsage / (1024 * 1024 * 1024)
	memoryCost := memoryGB * s.config.MemoryCostPerGB

	totalCost := cpuCost + memoryCost

	return &internal.NodeCost{
		Name:        nodeName,
		CPUCost:     cpuCost,
		MemoryCost:  memoryCost,
		TotalCost:   totalCost,
		CPUUsage:    cpuUsage,
		MemoryUsage: int64(memoryUsage),
		Timestamp:   time.Now(),
	}, nil
}
func (s *CostService) GetPodCosts(ctx context.Context, namespace string) ([]internal.PodCost, error) {
	pods, err := s.k8sClient.GetPods(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %v", err)
	}

	costs := make([]internal.PodCost, 0, len(pods.Items))
	for _, pod := range pods.Items {
		cost, err := s.calculatePodCost(ctx, pod.Namespace, pod.Name)
		if err != nil {
			continue
		}
		costs = append(costs, *cost)
	}

	return costs, nil
}

func (s *CostService) calculateNamespaceCost(ctx context.Context, namespace string) (*internal.NamespaceCost, error) {
	pods, err := s.k8sClient.GetPods(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods for namespace %s: %v", namespace, err)
	}
	var totalCPUCost, totalMemoryCost, totalStorageCost, totalNetworkCost float64

	podCosts := make([]internal.PodCost, 0, len(pods.Items))

	for _, pod := range pods.Items {
		podCost, err := s.calculatePodCost(ctx, pod.Namespace, pod.Name)
		if err != nil {
			continue
		}

		totalCPUCost += podCost.CPUCost
		totalMemoryCost += podCost.MemoryCost
		totalNetworkCost += podCost.NetworkCost
		totalStorageCost += podCost.StorageCost

		podCosts = append(podCosts, *podCost)
	}

	return &internal.NamespaceCost{
		Namespace:   namespace,
		CPUCost:     totalCPUCost,
		MemoryCost:  totalMemoryCost,
		StorageCost: totalStorageCost,
		NetworkCost: totalNetworkCost,
		TotalCost:   totalCPUCost + totalMemoryCost + totalStorageCost + totalNetworkCost,
		PodCount:    len(pods.Items),
		Pods:        podCosts,
		Timestamp:   time.Now(),
	}, nil

}

func (s *CostService) calculatePodCost(ctx context.Context, namespace, podName string) (*internal.PodCost, error) {
	cpuUsage, err := s.promClient.GetCPUUsage(ctx, namespace, podName)
	if err != nil {
		cpuUsage = 0
	}

	cpuCost := cpuUsage * s.config.CPUCostPerHour

	memoryUsage, err := s.promClient.GetMemoryUsage(ctx, namespace, podName)
	if err != nil {
		memoryUsage = 0
	}
	memoryGB := memoryUsage / (1024 * 1024 * 1024)
	memoryCost := memoryGB * s.config.MemoryCostPerGB

	rxBytes, txBytes, err := s.promClient.GetNetworkIO(ctx, namespace, podName)
	if err != nil {
		rxBytes = 0
		txBytes = 0
	}

	networkCost := (rxBytes + txBytes) * 0.000001

	storageCost := 0.01

	return &internal.PodCost{
		Name:        podName,
		Namespace:   namespace,
		CPUCost:     cpuCost,
		MemoryCost:  memoryCost,
		StorageCost: storageCost,
		NetworkCost: networkCost,
		TotalCost:   cpuCost + memoryCost + storageCost + networkCost,
		CPUUsage:    cpuUsage,
		MemoryUsage: int64(memoryUsage),
		// Status:      "Running",
		// Timestamp:   time.Time{},
		// CreatedAt:   time.Time{},
	}, nil
}

func (s *CostService) convertToCostHistory(cpuResult, memResult prometheus.RangeQueryResult) []internal.CostHistoryPoint {
	points := make([]internal.CostHistoryPoint, 0)

	// Assuming both CPU and memory results have the same timestamps
	timestampMap := make(map[int64]*internal.CostHistoryPoint)

	for _, series := range cpuResult.Data.Result {
		for _, value := range series.Values {
			timestamp := value.Timestamp
			cpuUsage := value.Value
			cpuCost := cpuUsage * s.config.CPUCostPerHour

			point, exists := timestampMap[timestamp.Unix()]
			if !exists {
				point = &internal.CostHistoryPoint{Timestamp: time.Unix(timestamp.Unix(), 0)}
				timestampMap[timestamp.Unix()] = point
			}
			point.CPUCost += cpuCost
		}
	}

	// Process Memory data
	for _, series := range memResult.Data.Result {
		for _, value := range series.Values {
			timestamp := value.Timestamp
			memoryUsage := value.Value
			memoryGB := memoryUsage / (1024 * 1024 * 1024)
			memoryCost := memoryGB * s.config.MemoryCostPerGB

			point, exists := timestampMap[timestamp.Unix()]
			if !exists {
				point = &internal.CostHistoryPoint{Timestamp: time.Unix(timestamp.Unix(), 0)}
				timestampMap[timestamp.Unix()] = point
			}
			point.MemoryCost += memoryCost
		}
	}

	// Finalize points
	for _, point := range timestampMap {
		point.TotalCost = point.CPUCost + point.MemoryCost
		points = append(points, *point)
	}

	// Optional: sort by timestamp
	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})

	return points
}
