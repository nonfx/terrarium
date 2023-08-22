package dbhelper

import (
	"log"

	"gorm.io/gorm/logger"
)

type connConfig struct {
	maxRetries       int
	retryIntervalSec int
	jitterLimitSec   int
	logger           logger.Interface
}

func getDefaultConfig() connConfig {
	return connConfig{
		maxRetries:       0,
		retryIntervalSec: 0,
		jitterLimitSec:   0,
		logger: logger.New(log.Default(), logger.Config{
			LogLevel: logger.Warn,
		}),
	}
}

type ConnOption func(*connConfig)

func WithRetries(maxRetries, retryIntervalSec, jitterLimitSec int) ConnOption {
	return func(cc *connConfig) {
		cc.maxRetries = maxRetries
		cc.retryIntervalSec = retryIntervalSec
		cc.jitterLimitSec = jitterLimitSec
	}
}

func WithLogger(writer logger.Writer, config logger.Config) ConnOption {
	return func(cc *connConfig) {
		cc.logger = logger.New(writer, config)
	}
}
