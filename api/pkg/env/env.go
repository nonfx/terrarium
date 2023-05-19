package env

import (
	"os"
	"strconv"
)

var PREFIX = ""

// GetEnvString returns the value of the specified environment variable,
// or the default value if the environment variable is not set.
func GetEnvString(name, defaultVal string) string {
	val, isSet := os.LookupEnv(PREFIX + name)
	if !isSet {
		return defaultVal
	}
	return val
}

// GetEnvInt returns the integer value of the specified environment variable,
// or the default value if the environment variable is not set or cannot be parsed as an integer.
func GetEnvInt(name string, defaultVal int) int {
	valStr, isSet := os.LookupEnv(PREFIX + name)
	if !isSet {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetEnvBool returns the boolean value of the specified environment variable,
// or the default value if the environment variable is not set or cannot be parsed as a boolean.
func GetEnvBool(name string, defaultVal bool) bool {
	valStr, isSet := os.LookupEnv(PREFIX + name)
	if !isSet {
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetEnvInt64 returns the int64 value of the specified environment variable,
// or the default value if the environment variable is not set or cannot be parsed as an int64.
func GetEnvInt64(name string, defaultVal int64) int64 {
	valStr, isSet := os.LookupEnv(PREFIX + name)
	if !isSet {
		return defaultVal
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}
