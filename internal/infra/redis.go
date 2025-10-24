package infra

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/redis/go-redis/v9"
)

func NewRedis(logger *zerolog.Logger) (*redis.Client, error) {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	db := 0
	var err error
	if os.Getenv("REDIS_DB") != "" {
		db, err = strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			return nil, err
		}
	}

	rd := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	return rd, nil
}
