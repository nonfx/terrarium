// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"log"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/confighelper"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/dbhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestAutoMigrate(t *testing.T) {
	t.Run("postgres", func(t *testing.T) {
		g := newConPG(t)
		dbObj, err := db.AutoMigrate(g)
		assert.NoError(t, err)
		assert.NotNil(t, dbObj)
	})
}

// Helpers used in multiple tests

func newConPG(t *testing.T) *gorm.DB {
	t.Helper()

	confighelper.LoadDefaults(map[string]interface{}{
		"db": map[string]interface{}{
			"host":     "localhost",
			"user":     "postgres",
			"password": "",
			"name":     "cc_terrarium",
			"port":     5432,
			"ssl_mode": false,
		},
	}, "tr")

	db, err := dbhelper.ConnectPG(
		confighelper.MustGetString("db.host"),
		confighelper.MustGetString("db.user"),
		confighelper.MustGetString("db.password"),
		confighelper.MustGetString("db.name"),
		confighelper.MustGetInt("db.port"),
		confighelper.MustGetBool("db.ssl_mode"),
		dbhelper.WithLogger(log.Default(), logger.Config{LogLevel: logger.Error}),
	)
	require.NoError(t, err)

	return db
}

type testConnector func(t *testing.T) *gorm.DB

func newConSQLite(t *testing.T) *gorm.DB {
	t.Helper()

	confighelper.LoadDefaults(map[string]interface{}{
		"db": map[string]interface{}{"dsn": ":memory:"},
	}, "tr")

	db, err := dbhelper.ConnectSQLite(
		confighelper.MustGetString("db.dsn"),
		dbhelper.WithLogger(log.Default(), logger.Config{LogLevel: logger.Error}),
	)
	require.NoError(t, err)

	return db
}

func getConnectorMap() map[string]testConnector {
	return map[string]testConnector{
		"postgres": newConPG,
		"sqlite":   newConSQLite,
	}
}
