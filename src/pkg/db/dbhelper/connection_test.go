// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package dbhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigPostgresGetDSN(t *testing.T) {
	testCases := []struct {
		name     string
		host     string
		user     string
		password string
		dbName   string
		port     int
		sslMode  bool
		expected string
	}{
		{
			name:     "sslMode disabled",
			host:     "localhost",
			user:     "test",
			password: "test",
			dbName:   "test",
			port:     5432,
			sslMode:  false,
			expected: "host=localhost user=test password=test dbname=test port=5432 sslmode=disable",
		},
		{
			name:     "sslMode enabled",
			host:     "localhost",
			user:     "test",
			password: "test",
			dbName:   "test",
			port:     5432,
			sslMode:  true,
			expected: "host=localhost user=test password=test dbname=test port=5432 sslmode=enable",
		},
		{
			name:     "no password",
			host:     "localhost",
			user:     "test",
			password: "",
			dbName:   "test",
			port:     5432,
			sslMode:  true,
			expected: "host=localhost user=test  dbname=test port=5432 sslmode=enable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dsn := (ConfigPostgres{Host: tc.host, User: tc.user, Password: tc.password, DBName: tc.dbName, Port: tc.port, SslMode: tc.sslMode}).GetDSN()
			assert.Equal(t, tc.expected, dsn)
		})
	}
}

func TestDialectorSwitcherSwitch(t *testing.T) {
	ds := DialectorSwitcher{
		ConfigPostgres: ConfigPostgres{},
		ConfigSQLite:   ConfigSQLite{},
	}
	tests := []struct {
		name              string
		ds                DialectorSwitcher
		dbType            DBDriver
		wantDialectorName string
		wantError         bool
	}{
		{
			name:              "success get postgres",
			ds:                ds,
			dbType:            DBDriverPostgres,
			wantDialectorName: "postgres",
		},
		{
			name:              "success get sqlite",
			ds:                ds,
			dbType:            DBDriverSQLite,
			wantDialectorName: "sqlite",
		},
		{
			name:      "fail undefined type",
			ds:        ds,
			dbType:    DBDriverUndefined,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDialector, gotErr := tt.ds.Switch(tt.dbType)
			if tt.wantError {
				assert.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.wantDialectorName, gotDialector.Name())
		})
	}
}

func TestDBDriver(t *testing.T) {
	tests := []struct {
		name    string
		fromStr string
		fromT   DBDriver
		wantT   DBDriver
		wantStr string
	}{
		{
			name:    "pg str to type",
			fromStr: "postgres",
			wantT:   DBDriverPostgres,
		},
		{
			name:    "sqlite str to type",
			fromStr: "sqlite",
			wantT:   DBDriverSQLite,
		},
		{
			name:    "invalid str to type",
			fromStr: "undefined",
			wantT:   DBDriverUndefined,
		},
		{
			name:    "pg type to str",
			fromT:   DBDriverPostgres,
			wantStr: "postgres",
		},
		{
			name:    "sqlite type to str",
			fromT:   DBDriverSQLite,
			wantStr: "sqlite",
		},
		{
			name:    "invalid type to str",
			fromT:   DBDriverUndefined,
			wantStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT := DBDriverFromStr(tt.fromStr)
			assert.Equal(t, tt.wantT, gotT)

			gotStr := tt.fromT.String()
			assert.Equal(t, tt.wantStr, gotStr)
		})
	}
}
