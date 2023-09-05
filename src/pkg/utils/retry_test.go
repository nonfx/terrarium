// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
			err := Retry(tc.maxRetries, tc.retryInterval, tc.jitterLimit, tc.funcToRetry)
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
