package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// SetupCompress sets up the compression middleware for the Fiber app
func SetupCompress(app *fiber.App) {
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
}
