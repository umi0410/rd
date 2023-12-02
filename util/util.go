package util

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var (
	DefaultLogger = log.StandardLogger()
)

func GetLogger(ctx context.Context) *log.Logger {
	val := ctx.Value("logger")
	if val != nil {
		logger, ok := val.(*log.Logger)
		if ok {
			return logger
		}
	}

	return DefaultLogger
}

