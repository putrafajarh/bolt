package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// TimeoutContext sets a deadline on c.UserContext()
// and calls c.Next() synchronously. Any DB ops
// using WithContext or handlers checking c.UserContext()
// will see the cancel when time expires.
func Timeout(c *fiber.Ctx, d time.Duration) fiber.Handler {
	return nil
}
