package config

import (
	"github.com/cldcvr/terrarium/api/pkg/dbhelper"
	"github.com/cldcvr/terrarium/api/pkg/env"
	"gorm.io/gorm"
)

func init() {
	env.PREFIX = "TR_" // Sets the environment prefix for the entire service
}

// DBHost returns the database host
func DBHost() string {
	return env.GetEnvString("DB_HOST", "localhost")
}

// DBUser returns the database user
func DBUser() string {
	return env.GetEnvString("DB_USER", "postgres")
}

// DBPassword returns the database password
func DBPassword() string {
	return env.GetEnvString("DB_PASSWORD", "")
}

// DBName returns the database name
func DBName() string {
	return env.GetEnvString("DB_NAME", "cc_terrarium")
}

// DBPort returns the database port
func DBPort() int {
	return env.GetEnvInt("DB_PORT", 5432)
}

// DBSSLMode returns the database SSL mode.
func DBSSLMode() bool {
	return env.GetEnvBool("DB_SSL_MODE", false)
}

// ConnectFromEnv establishes a connection to the database using the connection parameters from the environment.
func ConnectFromEnv() (*gorm.DB, error) {
	return dbhelper.ConnectPG(
		DBHost(),
		DBUser(),
		DBPassword(),
		DBName(),
		DBPort(),
		DBSSLMode(),
	)
}
