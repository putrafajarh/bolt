package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/fiber/v2"
	"github.com/putrafajarh/bolt/controllers"
	"github.com/putrafajarh/bolt/internal/infra"
	"github.com/putrafajarh/bolt/internal/middlewares"
	"github.com/rs/zerolog"

	_ "github.com/joho/godotenv/autoload"

	"gorm.io/gorm"
)

func main() {

	logger := infra.NewLogger()

	db, err := infra.NewDB(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed create new db connection")
	}

	rd, err := infra.NewRedis(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed create new redis connection")
	}

	http := infra.NewHttpServer(logger)
	registerRoutes(http, db)

	go func() {
		port := os.Getenv("APP_PORT")
		if err := http.Listen(fmt.Sprintf(":%s", port)); err != nil {
			logger.Fatal().Err(err).Msgf("failed start http server on port %s", port)
		}
	}()

	gracefulShutdown(http, logger, 30*time.Second,
		// Cleanup db connection
		func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				logger.Error().Err(err).Msg("failed to close db connection")
				return err
			}
			if err := sqlDB.Close(); err != nil {
				logger.Error().Err(err).Msg("failed to close db connection")
				return err
			}
			logger.Info().Msg("db connection closed")
			return nil
		},
		// Cleanup redis connection
		func(ctx context.Context) error {
			if err := rd.Close(); err != nil {
				logger.Error().Err(err).Msg("failed to close redis connection")
				return err
			}
			logger.Info().Msg("redis connection closed")
			return nil
		},
	)
}

func registerRoutes(app *fiber.App, db *gorm.DB) {
	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 3,
		Timeout:          10 * time.Second,
		SuccessThreshold: 2,
	})

	v1 := app.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	v1.Get("/ping",
		circuitbreaker.Middleware(cb),
		middlewares.WithTrx(db),
		controllers.HandlePing,
	)
}

func gracefulShutdown(app *fiber.App, logger *zerolog.Logger, timeout time.Duration, ops ...func(context.Context) error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	logger.Info().Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to shutdown http server gracefully")
	}

	var wg sync.WaitGroup
	for _, op := range ops {
		if op == nil {
			continue
		}
		wg.Add(1)
		go func(op func(context.Context) error) {
			defer wg.Done()
			if err := op(ctx); err != nil {
				logger.Error().Err(err).Msg("failed to shutdown gracefully")
			}
		}(op)
	}
	wg.Wait()

	logger.Info().Msg("server shutdown gracefully stopped")
}
