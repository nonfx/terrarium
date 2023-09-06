// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

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
		maxRetries:       2,
		retryIntervalSec: 2,
		jitterLimitSec:   2,
		logger: logger.New(log.Default(), logger.Config{
			LogLevel: logger.Warn,
		}),
	}
}

type ConnOption func(*connConfig)

// WithRetries configures db connection retries.
// maxRetries represents number of times the connection request should be retried on failure.
// one attempt is done regardless. and then re-attempt is done based on this number. i.e. <1
// means total 1 attempt. >1 means n+1 attempt in total.
// retryIntervalSec represents wait time (in seconds) before each retry attempt.
// jitterLimitSec represents jitter time limit (in seconds), such that a random time
// interval between 0 and this number is selected and added to the retry interval on each retry
// attempt to avoid high retry traffic on exact same time from all servers.
func WithRetries(maxRetries, retryIntervalSec, jitterLimitSec int) ConnOption {
	return func(cc *connConfig) {
		cc.maxRetries = maxRetries
		cc.retryIntervalSec = retryIntervalSec
		cc.jitterLimitSec = jitterLimitSec
	}
}

// WithLogger configures the logger instance used by the database manager.
// writer instance allows setting up the log destination. By default, we set it to
// standard `log.Default()`.
// config is the gorm logger config that allows setting the log level, etc.
// by default, we set log level to Warning.
func WithLogger(writer logger.Writer, config logger.Config) ConnOption {
	return func(cc *connConfig) {
		cc.logger = logger.New(writer, config)
	}
}
