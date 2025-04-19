package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/putrafajarh/bolt/docs"
)

// SetupSwagger sets up the Swagger documentation for the Fiber app
func SetupSwagger(app *fiber.App) {
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Title = "Bolt API"
	docs.SwaggerInfo.Description = "API documentation for Bolt"
	docs.SwaggerInfo.Version = "0.1.0"

	app.Get("/v1/swagger/*", swagger.HandlerDefault)
}
