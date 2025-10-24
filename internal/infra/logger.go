package infra

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// logWriter := zerolog.ConsoleWriter{Out: os.Stderr} // Human friendly readable
	logWriter := os.Stderr
	logger := zerolog.New(logWriter).With().Timestamp().Logger()

	return &logger
}
