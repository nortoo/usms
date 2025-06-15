package log

import (
	"github.com/nortoo/logger"
	"go.uber.org/zap"
)

var l *zap.Logger

func GetLogger() *zap.Logger {
	return l
}

func InitLogger(config string) error {
	var err error
	l, err = logger.New(config)
	return err
}
