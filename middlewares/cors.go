package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS sets up the CORS middleware for the Fiber app
func SetupCORS(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Change this to specific origins in production
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))
}
