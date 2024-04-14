package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {

	var l *zap.Logger

	lgConf := zap.NewProductionConfig()
	lgConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := lgConf.Build()
	if err != nil {
		return nil, err
	}

	return l, nil
}
