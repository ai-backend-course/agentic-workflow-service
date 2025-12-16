package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"agentic-workflow-service/internal/httpapi"
)

func main() {
	_ = godotenv.Load()

	// Read environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create fiber app
	app := fiber.New()

	// Basic logging middleware
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("%s %s", c.Method(), c.Path())
		return c.Next()
	})

	// Register routes
	httpapi.RegisterRoutes(app)

	// Start server
	log.Printf("Agent Workflow Service running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
