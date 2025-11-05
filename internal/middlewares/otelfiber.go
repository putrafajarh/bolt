package middlewares

import (
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// SetupOtelFiber sets up the otel middleware for the Fiber app
func SetupOtelFiber(app *fiber.App) {
	app.Use(otelfiber.Middleware(
		otelfiber.WithNext(func(c *fiber.Ctx) bool {
			// Skip otel middleware for ping endpoint
			if c.Path() == "/v1/ping" {
				return false
			}
			return false
		}),
	))
}

func GetTracer(c *fiber.Ctx) oteltrace.Tracer {
	return otel.GetTracerProvider().Tracer("gofiber-contrib-tracer-fiber")
}
