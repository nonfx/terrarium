package dbhelper

import (
	"fmt"

	"github.com/cldcvr/terrarium/src/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:generate mockery --name Dialector --srcpkg gorm.io/gorm

// Connect establishes a connection to the database using the given dialector and connection parameters.
func Connect(dialector gorm.Dialector, options ...ConnOption) (*gorm.DB, error) {
	var db *gorm.DB

	cc := getDefaultConfig()

	for _, op := range options {
		op(&cc)
	}

	err := utils.Retry(cc.maxRetries, cc.retryIntervalSec, cc.jitterLimitSec, func() error {
		var err error
		db, err = gorm.Open(dialector, &gorm.Config{
			Logger: cc.logger,
		})
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

// ConnectPG establishes a connection to a postgres database using the provided connection parameters.
func ConnectPG(host, user, password, dbName string, port int, sslMode bool, options ...ConnOption) (*gorm.DB, error) {
	dsn := createDSN(host, user, password, dbName, port, sslMode)
	return Connect(postgres.Open(dsn), options...)
}

func ConnectSQLite(dsn string, options ...ConnOption) (*gorm.DB, error) {
	return Connect(sqlite.Open(dsn), options...)
}
