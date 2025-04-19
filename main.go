package main

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/putrafajarh/bolt/middlewares"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		AppName:     "Bolt v0.1.0",
	})

	// Register Middlewares
	middlewares.SetupCORS(app)
	middlewares.SetupRequestID(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
