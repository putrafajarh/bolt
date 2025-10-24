package middlewares

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var (
	ErrNoLogger = fiber.NewError(fiber.StatusInternalServerError, "no logger found")
)

const (
	loggerKey = "logger"
)

func SetupLogger(app *fiber.App, logger *zerolog.Logger) {
	app.Use(WithRequestId(logger))
}

func WithRequestId(baseLogger *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var logger zerolog.Logger
		rid, ok := c.Locals("requestid").(string)
		if ok {
			logger = baseLogger.With().Str("request_id", rid).Logger()
		}

		// logger.UpdateContext(func(zc zerolog.Context) zerolog.Context {
		// 	rid := c.Locals("requestid").(string)
		// 	return zc.Str("request_id", rid)
		// })

		// logger.UpdateContext(func(zc zerolog.Context) zerolog.Context {
		// 	rid := c.Locals("requestid").(string)
		// 	return zc.Str("request_id", rid)
		// })

		c.Locals(loggerKey, &logger)

		h := fiberzerolog.New(fiberzerolog.Config{
			Logger: &logger,
		})

		return h(c)
	}
}

func GetLogger(c *fiber.Ctx) (*zerolog.Logger, error) {
	if v := c.Locals(loggerKey); v != nil {
		if tx, ok := v.(*zerolog.Logger); ok {
			return tx, nil
		}
	}

	return nil, ErrNoLogger
}
