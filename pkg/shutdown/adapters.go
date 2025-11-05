package shutdown

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
)

// HTTPShutdownAdapter adapts *http.Server to the Shutdowner interface
// Useful for Gin, Chi, and any std net/http server
// Usage: &HTTPShutdownAdapter{Server: srv}
type HTTPShutdownAdapter struct {
	Server *http.Server
}

func (s *HTTPShutdownAdapter) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

type FiberShutdownAdapter struct {
	App *fiber.App
}

func (f *FiberShutdownAdapter) Shutdown(ctx context.Context) error {
	return f.App.ShutdownWithContext(ctx)
}

// EchoShutdownAdapter adapts echo.Echo to the Shutdowner interface
type EchoShutdownAdapter struct {
	Server *echo.Echo
}

func (e *EchoShutdownAdapter) Shutdown(ctx context.Context) error {
	return e.Server.Shutdown(ctx)
}
