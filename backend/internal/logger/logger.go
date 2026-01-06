package logger

import (
	"go.uber.org/zap"
)

// New builds a zap logger with sensible defaults for structured logging.
func New(env string) (*zap.Logger, error) {
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
