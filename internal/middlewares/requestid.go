package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

var (
	HeaderXRequestID    = http.CanonicalHeaderKey(fiber.HeaderXRequestID)
	RequestIDContextKey = "requestid"
)

// SetupRequestID sets up the Request ID middleware for the Fiber app
func SetupRequestID(app *fiber.App) {
	app.Use(requestid.New(requestid.Config{
		Header:     HeaderXRequestID,
		ContextKey: RequestIDContextKey,
		Generator:  requestidGenerator,
	}))
}

func requestidGenerator() string {
	rid, err := uuid.NewV7()
	if err != nil {
		return utils.UUID()
	}
	return rid.String()
}

func GetRequestID(c *fiber.Ctx) string {
	return c.Locals(RequestIDContextKey).(string)
}
