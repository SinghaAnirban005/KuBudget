package models

type PodUtilization struct {
	Namespace     string  `json:"namespace"`
	PodName       string  `json:"podName"`
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
	Cost          float64 `json:"cost"`
}

type NodeUtilization struct {
	NodeName      string  `json:"nodeName"`
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
	Cost          float64 `json:"cost"`
	PodsRunning   int     `json:"podsRunning"`
}
