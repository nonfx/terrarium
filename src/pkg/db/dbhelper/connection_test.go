package dbhelper

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/cldcvr/terrarium/src/pkg/db/dbhelper/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestRetry(t *testing.T) {
	testCases := []struct {
		name           string
		maxRetries     int
		retryInterval  int
		jitterLimit    int
		funcToRetry    func() error
		expectedErrStr string
	}{
		{
			name:          "no error",
			maxRetries:    3,
			retryInterval: 1,
			jitterLimit:   1,
			funcToRetry: func() error {
				return nil
			},
		},
		{
			name:          "retry exceeds limit",
			maxRetries:    2,
			retryInterval: 1,
			jitterLimit:   1,
			funcToRetry: func() error {
				return errors.New("some error")
			},
			expectedErrStr: "failed after 2 retries: some error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			beginTime := time.Now()
			err := retry(tc.maxRetries, tc.retryInterval, tc.jitterLimit, tc.funcToRetry)
			duration := time.Since(beginTime)
			if tc.expectedErrStr != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrStr, err.Error())
				assert.GreaterOrEqual(t, duration, time.Duration(tc.retryInterval*(tc.maxRetries-1))*time.Second)
			} else {
				assert.NoError(t, err)
				assert.Less(t, duration, time.Second)
			}
		})
	}
}

func TestConnect(t *testing.T) {
	// Mock dialector
	mockErr := fmt.Errorf("mocked error")
	mockDialector := &mocks.Dialector{}
	mockDialector.On("Initialize", mock.Anything).Return(mockErr).Times(3)
	mockDialector.On("Initialize", mock.Anything).Return(nil)

	_, err := Connect(mockDialector, 2, 0, 0) // fail
	assert.ErrorIs(t, err, mockErr)

	db, err := Connect(mockDialector, 2, 0, 0) // succeed
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
