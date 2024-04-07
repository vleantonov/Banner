package logger

import (
	"fmt"
	"go.uber.org/zap"
)

const (
	envDev  = "development"
	envProd = "production"
)

func New(env string) (*zap.Logger, error) {

	var l *zap.Logger

	// TODO: add logger customization
	switch env {
	case envDev:
		l, _ = zap.NewDevelopment()
	case envProd:
		l, _ = zap.NewProduction()
	default:
		return nil, fmt.Errorf("unknown environment for logger init: %s", env)
	}

	return l, nil
}
