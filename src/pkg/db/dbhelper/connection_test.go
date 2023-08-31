// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package dbhelper

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db/dbhelper/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm/logger"
)

func TestCreateDSN(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dsn := createDSN(tc.host, tc.user, tc.password, tc.dbName, tc.port, tc.sslMode)
			assert.Equal(t, tc.expected, dsn)
		})
	}
}

func TestConnect(t *testing.T) {
	// Mock dialector
	mockErr := fmt.Errorf("mocked error")
	mockDialector := &mocks.Dialector{}
	mockDialector.On("Initialize", mock.Anything).Return(mockErr).Times(3)
	mockDialector.On("Initialize", mock.Anything).Return(nil).Once()

	logData := strings.Builder{}
	opts := []ConnOption{
		WithRetries(1, 0, 0),
		WithLogger(log.New(&logData, "", 0), logger.Config{LogLevel: logger.Info}),
	}

	_, err := Connect(mockDialector, opts...) // fail twice
	assert.ErrorIs(t, err, mockErr)

	db, err := Connect(mockDialector, opts...) // succeed on retry after failing once
	assert.NoError(t, err)
	assert.NotNil(t, db)

	assert.Equal(t, 3, strings.Count(logData.String(), "[error] failed to initialize database, got error mocked error"))
}
