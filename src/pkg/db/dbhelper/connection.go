// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package dbhelper

import (
	"fmt"

	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/rotisserie/eris"
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

	return db, eris.Wrap(err, "error connecting to database")
}

// createPostgresDSN creates the DSN (Data Source Name) string for the postgres database connection.
func createPostgresDSN(host, user, password, dbName string, port int, sslMode bool) string {
	sslModeStr := "disable"
	if sslMode {
		sslModeStr = "enable"
	}

	passwordStr := ""
	if password != "" {
		// omitting the password block when not set, allows the library to look at
		// other standard sources like `~/.pgpass`
		passwordStr = "password=" + password
	}

	return fmt.Sprintf("host=%s user=%s %s dbname=%s port=%d sslmode=%s", host, user, passwordStr, dbName, port, sslModeStr)
}

// ConnectPG establishes a connection to a postgres database using the provided connection parameters.
func ConnectPG(host, user, password, dbName string, port int, sslMode bool, options ...ConnOption) (*gorm.DB, error) {
	dsn := createPostgresDSN(host, user, password, dbName, port, sslMode)
	return Connect(postgres.Open(dsn), options...)
}

func ConnectSQLite(dsn string, options ...ConnOption) (*gorm.DB, error) {
	return Connect(sqlite.Open(dsn), options...)
}
