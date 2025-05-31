package handlers

import (
	"strconv"
	"time"

	"github.com/SinghaAnirban005/KuBudget/services"
	"github.com/gofiber/fiber/v2"
)

type CostHandler struct {
	costService *services.CostService
}

func NewCostHandler(costService *services.CostService) *CostHandler {
	return &CostHandler{
		costService: costService,
	}
}

func (h *CostHandler) GetCostOverview(c *fiber.Ctx) error {
	ctx := c.Context()

	overview, err := h.costService.GetCostOverview(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get cost overview",
			"details": err.Error(),
		})
	}

	return c.JSON(overview)
}

func (h *CostHandler) GetNamespaceCosts(c *fiber.Ctx) error {
	ctx := c.Context()

	costs, err := h.costService.GetNamespaceCosts(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get namespace costs",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"namespace_costs": costs,
		"count":           len(costs),
		"timestamp":       time.Now(),
	})

}

func (h *CostHandler) GetPodCosts(c *fiber.Ctx) error {
	ctx := c.Context()
	namespace := c.Query("namespace", "default")

	costs, err := h.costService.GetPodCosts(ctx, namespace)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get pod costs",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"pod_costs": costs,
		"namespace": namespace,
		"count":     len(costs),
		"timestamp": time.Now(),
	})
}

func (h *CostHandler) GetNodeCosts(c *fiber.Ctx) error {
	ctx := c.Context()

	costs, err := h.costService.GetNodeCosts(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get node costs",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"node_costs": costs,
		"count":      len(costs),
		"timestamp":  time.Now(),
	})
}

func (h *CostHandler) GetCostHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	hoursStr := c.Query("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid hours parameter",
		})
	}

	stepStr := c.Query("step", "1h")
	step, err := time.ParseDuration(stepStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid step parameter",
		})
	}

	namespace := c.Query("namespace", "")

	history, err := h.costService.GetCostHistory(ctx, time.Duration(hours)*time.Hour, step, namespace)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get cost history",
			"details": err.Error(),
		})
	}

	return c.JSON(history)
}
