package dbhelper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rotisserie/eris"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:generate mockery --name Dialector --srcpkg gorm.io/gorm

// Connect establishes a connection to the database using the given dialector and connection parameters.
func Connect(dialector gorm.Dialector, maxRetries, retryIntervalSec, jitterLimitSec int) (*gorm.DB, error) {
	var db *gorm.DB
	err := retry(maxRetries, retryIntervalSec, jitterLimitSec, func() error {
		var err error
		db, err = gorm.Open(dialector, &gorm.Config{})
		return err
	})

	return db, err
}

// createDSN creates the DSN (Data Source Name) string for the database connection.
func createDSN(host, user, password, dbName string, port int, sslMode bool) string {
	sslModeStr := "disable"
	if sslMode {
		sslModeStr = "enable"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, user, password, dbName, port, sslModeStr)
}

// retry executes the provided function with retry logic.
// in case of failures, the `fu` is retried `maxRetries` number of time,
// and wait for (`retryIntervalSec` + random `jitterLimitSec`) duration between each retry.
func retry(maxRetries int, retryIntervalSec, jitterLimitSec int, fu func() error) error {
	// Create a backoff function that adds jitter to the delay duration.
	backoff := func(attempt int) time.Duration {
		return (time.Second * time.Duration(retryIntervalSec)) + (time.Second * time.Duration(jitterLimitSec) * time.Duration(rand.Intn(100)) / 100)
	}

	var err error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			// Sleep before retrying
			time.Sleep(backoff(i))
		}
		err = fu()
		if err == nil {
			// Success, no need to retry
			return nil
		}
	}

	return eris.Wrapf(err, "failed after %d retries", maxRetries)
}

// ConnectPG establishes a connection to a postgres database using the provided connection parameters.
func ConnectPG(host, user, password, dbName string, port int, sslMode bool) (*gorm.DB, error) {
	dsn := createDSN(host, user, password, dbName, port, sslMode)
	return Connect(postgres.Open(dsn), 20, 3, 3)
}
