# Graceful Shutdown Usage

This package provides adapters for graceful shutdown with popular Go web frameworks.

## Fiber
```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/putrafajarh/bolt/pkg/shutdown"
    "github.com/rs/zerolog"
    "time"
)

app := fiber.New()
// ... setup routes ...

adapter := &shutdown.FiberShutdownAdapter{App: app}
shutdown.Graceful(adapter, logger, 30*time.Second, cleanupFuncs...)
```

## Echo
```go
import (
    "github.com/labstack/echo/v4"
    "github.com/putrafajarh/bolt/pkg/shutdown"
    "github.com/rs/zerolog"
    "time"
)

e := echo.New()
// ... setup routes ...

adapter := &shutdown.EchoShutdownAdapter{Server: e}
shutdown.Graceful(adapter, logger, 30*time.Second, cleanupFuncs...)
```

## Gin (and std net/http)
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/putrafajarh/bolt/pkg/shutdown"
    "github.com/rs/zerolog"
    "net/http"
    "time"
)

router := gin.New()
// ... setup routes ...

srv := &http.Server{
    Addr:    ":8080",
    Handler: router,
}
adapter := &shutdown.HTTPShutdownAdapter{Server: srv}
shutdown.Graceful(adapter, logger, 30*time.Second, cleanupFuncs...)
```

- `logger`: a zerolog logger instance
- `cleanupFuncs`: optional resource cleanup functions, signature: `func(ctx context.Context) error`
- `timeout`: shutdown timeout duration
