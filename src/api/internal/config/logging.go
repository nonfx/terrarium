// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"log"
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/confighelper"
	"github.com/sirupsen/logrus"
)

// LogFormat JSON or TEXT (case insensitive)
func LogFormat() string {
	return confighelper.MustGetString("log.format")
}

// LogLevel one of: "panic", "fatal", "error", "warn", "warning", "info", "debug", "trace" (case sensitive)
func LogLevel() string {
	return confighelper.MustGetString("log.level")
}

// LogPrettyPrint for JSON it means indentation on, for TEXT it means force color
func LogPrettyPrint() bool {
	return confighelper.MustGetBool("log.pretty_print")
}

// LoggerConfig sets up the given logger with defined configuration.
//
// Example:
// To update standard logger configuration, do this:
//
// LoggerConfig(logrus.StandardLogger())
func LoggerConfig(logger *logrus.Logger) {
	// Set log formatter
	formatter := LogFormat()
	if strings.EqualFold(formatter, "TEXT") {
		logger.SetFormatter(&logrus.TextFormatter{ForceColors: LogPrettyPrint()})
	} else if strings.EqualFold(formatter, "JSON") {
		logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: LogPrettyPrint()})
	}

	// Set log level
	levelStr := LogLevel()
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.Debugf("failed to parse log level string %s: %v", levelStr, err)
	} else {
		logger.SetLevel(level)
	}
}

// LoggerConfigDefault sets up the default loggers with defined configuration.
func LoggerConfigDefault() {
	logger := logrus.StandardLogger()
	LoggerConfig(logger) // setup `logger` object from config

	// update default standard logger
	log.Default().SetOutput(logger.WriterLevel(logrus.DebugLevel))
}
