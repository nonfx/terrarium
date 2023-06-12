package config

import (
	"strings"

	"github.com/cldcvr/terrarium/api/pkg/env"
	log "github.com/sirupsen/logrus"
)

// LogFormat JSON or TEXT (case insensitive)
func LogFormat() string {
	return env.GetEnvString("LOG_FORMAT", "JSON")
}

// LogLevel one of: "panic", "fatal", "error", "warn", "warning", "info", "debug", "trace" (case sensitive)
func LogLevel() string {
	return env.GetEnvString("LOG_LEVEL", "info")
}

// LogPrettyPrint for JSON it means indentation on, for TEXT it means force color
func LogPrettyPrint() bool {
	return env.GetEnvBool("LOG_PRETTY_PRINT", false)
}

// SetupLogger sets up teh given logger with defined configuration.
//
// Example:
// To update standard logger configuration, do this:
//
// SetupLogger(logrus.StandardLogger())
func SetupLogger(logger *log.Logger) {
	// Set log formatter
	formatter := LogFormat()
	if strings.EqualFold(formatter, "TEXT") {
		logger.SetFormatter(&log.TextFormatter{ForceColors: LogPrettyPrint()})
	} else if strings.EqualFold(formatter, "JSON") {
		logger.SetFormatter(&log.JSONFormatter{PrettyPrint: LogPrettyPrint()})
	}

	// Set log level
	levelStr := LogLevel()
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		log.Debugf("failed to parse log level string %s: %v", levelStr, err)
	} else {
		logger.SetLevel(level)
	}
}
