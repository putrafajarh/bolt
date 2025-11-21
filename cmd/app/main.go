package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub/v2"
	"github.com/putrafajarh/bolt/pkg/shutdown"
	"golang.org/x/sync/errgroup"

	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/putrafajarh/bolt/controllers"
	"github.com/putrafajarh/bolt/internal/consumer"
	ctrl "github.com/putrafajarh/bolt/internal/controller"
	"github.com/putrafajarh/bolt/internal/infra"
	"github.com/putrafajarh/bolt/internal/middlewares"
	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"

	_ "github.com/joho/godotenv/autoload"

	"gorm.io/gorm"
)

func main() {

	logger := infra.NewLogger()

	appName := os.Getenv("APP_NAME")
	tp, mp, err := infra.SetupOtel(appName)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to setup otel")
	}

	db, err := infra.NewDB(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed create new db connection")
	}

	rd, err := infra.NewRedis(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed create new redis connection")
	}

	fs, err := infra.NewFirestoreClient()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed create new firestore connection")
	}

	http := infra.NewHttpServer(logger)
	registerRoutes(http, db)

	go func() {
		port := os.Getenv("APP_PORT")
		if err := http.Listen(fmt.Sprintf(":%s", port)); err != nil {
			logger.Fatal().Err(err).Msgf("failed start http server on port %s", port)
		}
	}()

	adapter := shutdown.FiberShutdownAdapter{App: http}
	shutdown.Graceful(&adapter, logger, 30*time.Second,
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
		// Cleanup firestore connection
		func(ctx context.Context) error {
			if err := fs.Close(); err != nil {
				logger.Error().Err(err).Msg("failed to close firestore connection")
				return err
			}
			logger.Info().Msg("firestore connection closed")
			return nil
		},
		// Shutdown TraceProvider
		func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
		// Shutdown MeterProvider
		func(ctx context.Context) error {
			return mp.Shutdown(ctx)
		},
	)
}

func registerRoutes(app *fiber.App, db *gorm.DB) {
	if err := runtimemetrics.Start(); err != nil {
		panic(err)
	}

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	middlewares.SetupOtelFiber(app)
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

	jira := v1.Group("/jira", middlewares.WithTrx(db))
	jira.Get("/me", ctrl.Me)
	jira.Get("/issue", ctrl.GetIssue)
}

// Put all consumers here
func runConsumers(pubsubClient *pubsub.Client, firestoreClient *firestore.Client, db *gorm.DB) {
	var g errgroup.Group
	g.Go(func() error {
		return consumer.DailySyncIssues(pubsubClient, firestoreClient, db)
	})
	g.Wait()
}
