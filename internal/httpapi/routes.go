package httpapi

import (
	"github.com/gofiber/fiber/v2"
)

type Handlers struct{}

func RegisterRoutes(app *fiber.App) {
	h := Handlers{}

	app.Post("/run", h.runAgent)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

}
