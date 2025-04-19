package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// SetupRequestID sets up the Request ID middleware for the Fiber app
func SetupRequestID(app *fiber.App) {
	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
	}))
}
