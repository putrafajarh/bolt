package infra

import (
	"errors"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/putrafajarh/bolt/internal/middlewares"
	"github.com/rs/zerolog"
)

func NewHttpServer(logger *zerolog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		AppName:      "Bolt v0.1.0",
		ErrorHandler: fiberErrorHandler,
	})

	registerGlobalMiddlewares(app, logger)

	return app
}

func fiberErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func registerGlobalMiddlewares(app *fiber.App, logger *zerolog.Logger) {
	middlewares.SetupRequestID(app)
	middlewares.SetupLogger(app, logger)
	middlewares.SetupCORS(app)
	middlewares.SetupCompress(app)
	middlewares.SetupSwagger(app)
}
