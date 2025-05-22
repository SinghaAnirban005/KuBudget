package internal

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	PrometheusURL    string
	KubeConfigPath   string
	Port             string
	MetricsInterval  time.Duration
	CPUCostPerHour   float64
	MemoryCostPerGB  float64
	StorageCostPerGB float64
}

func LoadConfig() *Config {
	return &Config{
		PrometheusURL:    getEnv("PROMETHEUS_URL", "http://localhost:9090"),
		KubeConfigPath:   getEnv("KUBE_CONFIG_PATH", ""),
		Port:             getEnv("PORT", "8080"),
		MetricsInterval:  getDurationEnv("METRICS_INTERVAL", 30*time.Second),
		CPUCostPerHour:   getFloatEnv("CPU_COST_PER_HOUR", 0.048),
		MemoryCostPerGB:  getFloatEnv("MEMORY_COST_PER_GB", 0.0067),
		StorageCostPerGB: getFloatEnv("STORAGE_COST_PER_GB", 0.00014),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}
