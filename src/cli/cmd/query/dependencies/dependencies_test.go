// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package dependencies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/stretchr/testify/assert"
)

func Test_fectchDependencies(t *testing.T) {
	config.LoadDefaults()

	mockDependencyID := "mock-dep-id"
	mockDB := &mocks.DB{}

	mockDB.On("FetchDependencyByInterfaceID", mockDependencyID).
		Return(nil, fmt.Errorf("mock error")).Once()

	mockDependency := &terrariumpb.Dependency{
		InterfaceId: mockDependencyID,
		Title:       "Mock Title",
		Description: "Mock Description",
		Inputs:      "Mock Inputs",
		Outputs:     "Mock Outputs",
	}
	mockDB.On("FetchDependencyByInterfaceID", mockDependencyID).Return(mockDependency, nil).Times(2)

	config.SetDBMocks(mockDB)

	expectedTableOutput := `+--------------+------------+------------------+-------------+--------------+
| INTERFACE ID |   TITLE    |   DESCRIPTION    |   INPUTS    |   OUTPUTS    |
+--------------+------------+------------------+-------------+--------------+
| mock-dep-id  | Mock Title | Mock Description | Mock Inputs | Mock Outputs |
+--------------+------------+------------------+-------------+--------------+
`

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		expectedOutput string
	}{
		{
			name:    "db failure",
			args:    []string{"-i", mockDependencyID},
			wantErr: true,
		},
		{
			name:           "fetch dependencies in json format",
			args:           []string{"-o", "json", "-i", mockDependencyID},
			expectedOutput: `{"interfaceId":"mock-dep-id","title":"Mock Title","description":"Mock Description","inputs":"Mock Inputs","outputs":"Mock Outputs"}`,
		}, {
			name:           "fetch dependencies in table format",
			args:           []string{"-o", "table", "-i", mockDependencyID},
			expectedOutput: expectedTableOutput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd()
			cmd.SetArgs(tt.args)
			buffer := &bytes.Buffer{}
			cmd.SetOutput(buffer)
			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, buffer.String())
			}
		})
	}
}
