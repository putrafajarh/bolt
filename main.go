package main

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/putrafajarh/bolt/controllers"
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
	middlewares.SetupCompress(app)
	middlewares.SetupSwagger(app)

	v1 := app.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	v1.Get("/ping", controllers.HandlePing)

	log.Fatal(app.Listen(":8000"))
}
