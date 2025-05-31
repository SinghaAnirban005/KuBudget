package main

import (
	"log"
	"os"

	"github.com/SinghaAnirban005/KuBudget/handlers"
	"github.com/SinghaAnirban005/KuBudget/internal"
	"github.com/SinghaAnirban005/KuBudget/pkg/kubernetes"
	"github.com/SinghaAnirban005/KuBudget/pkg/prometheus"
	"github.com/SinghaAnirban005/KuBudget/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	cfg := internal.LoadConfig()

	k8sClient, err := kubernetes.NewClient(cfg.KubeConfigPath)
	if err != nil {
		log.Fatalf("Failed to create kubernetes cluster,  %v", err)
	}

	promClient := prometheus.NewClient(cfg.PrometheusURL)

	costService := services.NewCostService(k8sClient, promClient)
	metricsService := services.NewMetricsService(k8sClient, promClient)

	costHandler := handlers.NewCostHandler(costService)
	healthHandler := handlers.NewHealthHandler()
	metricsHandler := handlers.NewMetricsHandler(metricsService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	api := app.Group("/api/v1")

	api.Get("/health", healthHandler.Health)
	api.Get("/ready", healthHandler.Ready)
	api.Get("/costs/overview", costHandler.GetCostOverview)
	api.Get("/costs/namespaces", costHandler.GetNamespaceCosts)
	api.Get("/costs/pods", costHandler.GetPodCosts)
	api.Get("/costs/nodes", costHandler.GetNodeCosts)
	api.Get("/costs/history", costHandler.GetCostHistory)

	api.Get("/metrics/prometheus", metricsHandler.GetPrometheusMetrics)
	api.Get("/metrics/cluster", metricsHandler.GetClusterMetrics)
	api.Get("/metrics/resource-usage", metricsHandler.GetResourceUsage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server is starting on port %s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}
