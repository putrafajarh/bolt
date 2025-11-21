package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	ErrNoTx = fiber.NewError(fiber.StatusInternalServerError, "no transaction found")
)

const (
	key = "db_tx"
)

var skipMethods = map[string]struct{}{
	fiber.MethodOptions: {},
	// fiber.MethodGet:     {},
	fiber.MethodConnect: {},
	fiber.MethodHead:    {},
}

func WithTrx(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		if _, ok := skipMethods[c.Method()]; ok {
			return c.Next()
		}

		logger, err := GetLogger(c)
		if err != nil {
			return err
		}
		tx := db.WithContext(logger.WithContext(c.UserContext())).Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			}
		}()

		c.Locals(key, tx)

		if err := c.Next(); err != nil {
			logger.Error().Err(err).Msg("Rollback transaction")
			tx.Rollback()
			return err
		}

		if isOk(c.Response().StatusCode()) {
			logger.Debug().Msg("Commit transaction")
			return tx.Commit().Error
		}

		logger.Debug().Msg("Rollback transaction")
		tx.Rollback()
		return nil
	}
}

// IsOK returns true if the status code is 2xx
func isOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// GORM database transaction middleware
func SetupTransactionMiddleware(app *fiber.App, db *gorm.DB) {
	app.Use(WithTrx(db))
}

func GetTx(c *fiber.Ctx) (*gorm.DB, error) {
	if v := c.Locals(key); v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx, nil
		}
	}

	return nil, ErrNoTx
}
