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
