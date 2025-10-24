package gormlogger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	gormLogger "gorm.io/gorm/logger"
)

const (
	DefaultLogLevel                  = gormLogger.Info
	DefaultSlowThreshold             = 200 * time.Millisecond
	DefaultIgnoreRecordNotFoundError = false
	DefaultSkipCaller                = 2
	DefaultParameterizedQueries      = false
)

type Logger struct {
	Logger                    *zerolog.Logger
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  gormLogger.LogLevel
	ParameterizedQueries      bool
	SkipCaller                int
}

var CALLER_PREFIX_PATH = ""

type Option func(*Logger)

func WithLogLevel(level gormLogger.LogLevel) Option {
	return func(l *Logger) {
		l.LogLevel = level
	}
}

func WithSlowThreshold(threshold time.Duration) Option {
	return func(l *Logger) {
		l.SlowThreshold = threshold
	}
}

func WithIgnoreRecordNotFoundError(ignore bool) Option {
	return func(l *Logger) {
		l.IgnoreRecordNotFoundError = ignore
	}
}

func WithParameterizedQueries(parameterized bool) Option {
	return func(l *Logger) {
		l.ParameterizedQueries = parameterized
	}
}

func WithSkipCaller(skip int) Option {
	return func(l *Logger) {
		l.SkipCaller = skip
	}
}

func setCallerPrefixPath(logger *zerolog.Logger) {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed get current working directory")
	}
	CALLER_PREFIX_PATH = cwd
}

func New(logger *zerolog.Logger, opts ...Option) *Logger {
	setCallerPrefixPath(logger)

	l := &Logger{
		Logger:                    logger,
		LogLevel:                  DefaultLogLevel,
		SlowThreshold:             DefaultSlowThreshold,
		IgnoreRecordNotFoundError: DefaultIgnoreRecordNotFoundError,
		SkipCaller:                DefaultSkipCaller,
		ParameterizedQueries:      DefaultParameterizedQueries,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.LogLevel = level
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= gormLogger.Info {
		l.logger(ctx).Info().Msgf(msg, data...)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= gormLogger.Warn {
		l.logger(ctx).Warn().Msgf(msg, data...)
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= gormLogger.Error {
		l.logger(ctx).Error().Msgf(msg, data...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	rowsAffected := strconv.FormatInt(rows, 10)
	if rows == -1 {
		rowsAffected = "-"
	}

	logger := l.logger(ctx)
	duration := float64(elapsed.Nanoseconds()) / 1e6 // in milliseconds

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		logger.Error().Err(err).Str("caller", fileWithLineNum()).Str("sql", sql).Str("rows", rowsAffected).Float64("duration", duration).Send()
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		msg := fmt.Sprintf("Slow SQL >= %v", l.SlowThreshold)
		logger.Warn().Str("caller", fileWithLineNum()).Str("sql", sql).Str("rows", rowsAffected).Float64("duration", duration).Msg(msg)
	case l.LogLevel == gormLogger.Info:
		logger.Info().Str("caller", fileWithLineNum()).Str("sql", sql).Str("rows", rowsAffected).Float64("duration", duration).Send()
	}
}

func (l *Logger) logger(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// func (l *Logger) addFields(ev *zerolog.Event, data ...any) *zerolog.Event {
// 	for i := 0; i < len(data); i += 2 {
// 		key, ok := data[i].(string)
// 		if !ok {
// 			continue
// 		}
// 		ev = ev.Interface(key, data[i+1])
// 	}
// 	return ev
// }

func (l *Logger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if l.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
)

// fileWithLineNum return the file name and line number of the current file
func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		file = strings.TrimPrefix(file, CALLER_PREFIX_PATH)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		default:
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}
