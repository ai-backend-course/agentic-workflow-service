package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	app := fiber.New()

	// Simple request logging middleware
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("%s %s", c.Method(), c.Path())
		return c.Next()
	})

	return app
}
