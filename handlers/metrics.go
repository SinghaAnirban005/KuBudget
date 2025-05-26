package handlers

import (
	"time"

	"github.com/SinghaAnirban005/KuBudget/services"
	"github.com/gofiber/fiber/v2"
)

type MetricsHandler struct {
	metricsService *services.MetricsService
}

func NewMetricsHandler(metricsService *services.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
}

func (h *MetricsHandler) GetPrometheusMetrics(c *fiber.Ctx) error {
	ctx := c.Context()

	namespace := c.Query("namespace", "")
	pod := c.Query("pod", "")

	metrics, err := h.metricsService.GetPrometheusMetrics(ctx, namespace, pod)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get Prometheus metrics",
			"details": err.Error(),
		})
	}

	return c.JSON(metrics)
}

func (h *MetricsHandler) GetClusterMetrics(c *fiber.Ctx) error {
	ctx := c.Context()

	metrics, err := h.metricsService.GetClusterMetrics(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get cluster metrics",
			"details": err.Error(),
		})
	}

	return c.JSON(metrics)
}

func (h *MetricsHandler) GetResourceUsage(c *fiber.Ctx) error {
	ctx := c.Context()

	namespace := c.Query("namespace", "")
	pod := c.Query("pod", "")

	usage, err := h.metricsService.GetResourceUsage(ctx, namespace, pod)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get resource usage",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"resource_usage": usage,
		"namespace":      namespace,
		"pod":            pod,
		"timestamp":      time.Now(),
	})
}
