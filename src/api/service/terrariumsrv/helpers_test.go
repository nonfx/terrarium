// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"
	"os"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase[REQ, RESP any] struct {
	name      string
	preCall   func(*testing.T, TestCase[REQ, RESP])
	req       *REQ
	wantResp  *RESP
	wantErr   interface{}
	postCheck func(*testing.T, TestCase[REQ, RESP], *RESP)
	mockDB    *mocks.DB // INTERNAL. DO NOT SET THIS. Add mocks to this
}

type TestCases[REQ, RESP any] []TestCase[REQ, RESP]

type serviceFunc[REQ, RESP any] func(*Service) func(context.Context, *REQ) (*RESP, error)

func (tc TestCase[REQ, RESP]) Run(t *testing.T, fu serviceFunc[REQ, RESP]) {
	// Run the test cases
	t.Run(tc.name, func(t *testing.T) {
		// Set up the mock behavior
		tc.mockDB = &mocks.DB{}
		service := New(tc.mockDB)

		if tc.preCall != nil {
			tc.preCall(t, tc)
		}

		resp, err := fu(service)(context.Background(), tc.req)

		// Assert the results
		if tc.wantErr == nil {
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
		} else {
			assert.Equal(t, tc.wantErr, err, "\nWant Err: %s\n Got Err: %s", tc.wantErr, err)
		}

		// Verify that the expected methods were called on the mock
		tc.mockDB.AssertExpectations(t)

		if tc.postCheck != nil {
			tc.postCheck(t, tc, resp)
		}
	})
}

func (testCases TestCases[REQ, RESP]) Run(t *testing.T, fu serviceFunc[REQ, RESP]) {
	for _, tc := range testCases {
		tc.Run(t, fu)
	}
}

func writeToFile(t *testing.T, filename, content string) {
	file, err := os.Create(filename)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(content)
	require.NoError(t, err)
}
