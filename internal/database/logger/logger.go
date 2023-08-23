package logger

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	nativeLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type logger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Logger                zerolog.Logger
}

func New() *logger {
	return &logger{
		Logger:                log.Logger,
		SkipErrRecordNotFound: false,
	}
}

func (l *logger) LogMode(nativeLogger.LogLevel) nativeLogger.Interface {
	return l
}

func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Info().Str("tag", core.DatabaseLogTag).Msgf(s, args)
}

func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Warn().Str("tag", core.DatabaseLogTag).Msgf(s, args)
}

func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Error().Str("tag", core.DatabaseLogTag).Msgf(s, args)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := map[string]interface{}{
		"sql":     sql,
		"latency": fmt.Sprintf("%v", elapsed),
		"rows":    rows,
	}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		l.Logger.Error().Err(err).Fields(fields).Str("tag", core.DatabaseLogTag).Msg("query error")
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.Warn().Fields(fields).Str("tag", core.DatabaseLogTag).Msg("slow query")
		return
	}

	l.Logger.Debug().Fields(fields).Str("tag", core.DatabaseLogTag).Msg("query")
}
