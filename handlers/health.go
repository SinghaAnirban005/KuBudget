package handlers

import (
	"time"

	"github.com/SinghaAnirban005/KuBudget/internal"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	healthStatus := internal.HealthStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Checks: map[string]string{
			"server": "ok",
			"api":    "ok",
		},
	}

	return c.JSON(healthStatus)
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	readinessStatus := internal.HealthStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Checks: map[string]string{
			"kubernetes": "ok",
			"prometheus": "ok",
			"api":        "ready",
		},
	}

	return c.JSON(readinessStatus)
}
