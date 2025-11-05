package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

// GracefulShutdown handles OS signals and gracefully shuts down the server using the provided shutdown function(s).
func Graceful(
	server Shutdowner,
	logger *zerolog.Logger,
	timeout time.Duration,
	cleanupFuncs ...func(context.Context) error,
) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info().Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to shutdown server gracefully")
	}

	var wg sync.WaitGroup
	for _, fn := range cleanupFuncs {
		if fn == nil {
			continue
		}
		wg.Add(1)
		go func(f func(context.Context) error) {
			defer wg.Done()
			if err := f(ctx); err != nil {
				logger.Error().Err(err).Msg("failed to shutdown gracefully")
			}
		}(fn)
	}
	wg.Wait()

	logger.Info().Msg("server shutdown gracefully stopped")
}
