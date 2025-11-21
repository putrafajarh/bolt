package infra

import (
	"fmt"
	"os"

	repo "github.com/putrafajarh/bolt/internal/repository/postgres"
	"github.com/putrafajarh/bolt/pkg/gormlogger"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(logger *zerolog.Logger) (*gorm.DB, error) {

	gormLogger := gormlogger.New(logger,
		gormlogger.WithParameterizedQueries(true),
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&repo.User{}, &repo.Project{}, &repo.Issue{})

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
