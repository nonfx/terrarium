// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/cldcvr/terrarium/src/pkg/confighelper"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/dbhelper"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

var mockdb *mocks.DB

// DBHost returns the database host
func DBHost() string {
	return confighelper.MustGetString("db.host")
}

// DBUser returns the database user
func DBUser() string {
	return confighelper.MustGetString("db.user")
}

// DBPassword returns the database password
func DBPassword() string {
	return confighelper.MustGetString("db.password")
}

// DBName returns the database name
func DBName() string {
	return confighelper.MustGetString("db.name")
}

// DBPort returns the database port
func DBPort() int {
	return confighelper.MustGetInt("db.port")
}

// DBSSLMode returns the database SSL mode.
func DBSSLMode() bool {
	return confighelper.MustGetBool("db.ssl_mode")
}

// DBType returns the database type chosen. Default is postgres
func DBType() string {
	return confighelper.MustGetString("db.type")
}

// DB_DSN returns the SQLite Data Source Name. Default is "embedded.db"
func DBDSN() string {
	return confighelper.MustGetString("db.dsn")
}

// DBRetryAttempts returns the max number of re-ties on DB connection request failure
func DBRetryAttempts() int {
	return confighelper.MustGetInt("db.retry_attempts")
}

// DBRetryInterval returns interval to wait for before retrying. (in seconds)
func DBRetryInterval() int {
	return confighelper.MustGetInt("db.retry_interval_sec")
}

// DBRetryInterval returns additional random time (in seconds) between 0 and this, to be added to the wait time.
func DBRetryJitter() int {
	return confighelper.MustGetInt("db.retry_jitter_sec")
}

// DBConnect establishes a connection to the database using the connection parameters from the environment.
func DBConnect() (db.DB, error) {
	var g *gorm.DB
	var err error
	switch DBType() {
	case "mock":
		if mockdb == nil {
			return nil, eris.New("mocked err: connection failed")
		}
		return mockdb, nil
	case "sqlite":
		g, err = dbhelper.ConnectSQLite(DBDSN())
		if err != nil {
			return nil, eris.Wrap(err, "could not establish a connection to the database")
		}
	default:
		g, err = dbhelper.ConnectPG(
			DBHost(),
			DBUser(),
			DBPassword(),
			DBName(),
			DBPort(),
			DBSSLMode(),
			dbhelper.WithRetries(DBRetryAttempts(), DBRetryInterval(), DBRetryJitter()),
		)
		if err != nil {
			return nil, eris.Wrap(err, "could not establish a connection to the database")
		}
	}
	return db.AutoMigrate(g)
}
