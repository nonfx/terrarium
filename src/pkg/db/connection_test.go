// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/confighelper"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/dbhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAutoMigrate(t *testing.T) {
	t.Run("postgres", func(t *testing.T) {
		g := newCon(dbhelper.DBDriverPostgres)(t)
		dbObj, err := db.AutoMigrate(g)
		assert.NoError(t, err)
		assert.NotNil(t, dbObj)
	})
}

// Helpers used in multiple tests

type testConnector func(t *testing.T) *gorm.DB

func newCon(dbType dbhelper.DBDriver) testConnector {
	return func(t *testing.T) *gorm.DB {

		t.Helper()

		confighelper.LoadDefaults(map[string]interface{}{
			"db": map[string]interface{}{
				"host":     "localhost",
				"user":     "postgres",
				"password": "",
				"name":     "cc_terrarium",
				"port":     5432,
				"ssl_mode": false,
				"dsn":      ":memory:",
			},
		}, "tr")

		config := dbhelper.DialectorSwitcher{
			ConfigPostgres: dbhelper.ConfigPostgres{
				Host:     confighelper.MustGetString("db.host"),
				User:     confighelper.MustGetString("db.user"),
				Password: confighelper.MustGetString("db.password"),
				DBName:   confighelper.MustGetString("db.name"),
				Port:     confighelper.MustGetInt("db.port"),
				SslMode:  confighelper.MustGetBool("db.ssl_mode"),
			},
			ConfigSQLite: dbhelper.ConfigSQLite{
				DSN: confighelper.MustGetString("db.dsn"),
			},
		}

		db, err := config.Connect(dbType)
		require.NoError(t, err)

		return db
	}
}

func getConnectorMap() map[string]testConnector {
	return map[string]testConnector{
		"postgres": newCon(dbhelper.DBDriverPostgres),
		"sqlite":   newCon(dbhelper.DBDriverSQLite),
	}
}
