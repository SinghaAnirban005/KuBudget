package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New(fiber.Config{
		AppName:       "K8s cost analyzer",
		CaseSensitive: true,
		StrictRouting: true,
	})

	// app.Use(app.Server().Logger.New())
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck)
	api.Get("/cluster-cost")
	api.Get("/pod-utilization")
	api.Get("/node-utilization")
	api.Get("/recommendations")

	app.Listen(":4000")
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
