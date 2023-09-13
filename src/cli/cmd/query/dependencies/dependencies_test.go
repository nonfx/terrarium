// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package dependencies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_fetchDependencies(t *testing.T) {
// 	// Setup and mocks
// 	config.LoadDefaults()
// 	mockDB := &mocks.DB{}
// 	mockUuid1 := uuid.New()

// 	// The exact number of FilterOptions and their order matters.
// 	// Ensure you are mocking it correctly.
// 	mockDependencyData := db.Dependencies{
// 		{
// 			Model:       db.Model{ID: mockUuid1},
// 			InterfaceID: "mockInterface123",
// 			Title:       "SampleDependency",
// 			Description: "This is a mock dependency",
// 			Inputs:      nil, // or an appropriate `jsonschema.Node` mock
// 			Outputs:     nil, // or an appropriate `jsonschema.Node` mock
// 			ExtendsID:   "mockExtend456",
// 			Taxonomy:    nil, // or an appropriate `Taxonomy` mock
// 		},
// 	}
// 	mockDB.On("FetchAllDependency",
// 		mock.AnythingOfType("db.FilterOption"),
// 		mock.AnythingOfType("db.FilterOption"),
// 	).Return(mockDependencyData, nil)

// 	config.SetDBMocks(mockDB)

// 	expectedJSON := fmt.Sprintf(`{
// 		"dependencies": [{
// 			"id": "%s",
// 			"interfaceId": "mockInterface123",
// 			"title": "SampleDependency",
// 			"description": "This is a mock dependency",
// 			"inputs": null,
// 			"outputs": null,
// 			"extendsId": "mockExtend456",
// 			"taxonomy": null
// 		}]
// 	}`, mockUuid1.String())

// 	tests := []struct {
// 		Name           string
// 		Args           []string
// 		WantErr        bool
// 		ExpectedOutput string
// 	}{
// 		{
// 			Name:           "fetch dependencies in json format",
// 			Args:           []string{"-o", "json"},
// 			ExpectedOutput: expectedJSON,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			cmd := &cobra.Command{}
// 			cmd.SetArgs(tt.Args)

// 			err := fetchDependencies(cmd, tt.Args)
// 			if (err != nil) != tt.WantErr {
// 				t.Errorf("fetchDependencies() error = %v, wantErr %v", err, tt.WantErr)
// 				return
// 			}
// 		})
// 	}
// }

func TestNewCmd(t *testing.T) {
	cmd := NewCmd()

	assert.Equal(t, "dependencies", cmd.Use)
	assert.Equal(t, "List dependency details matching the dependency interface id", cmd.Short)

	searchFlag := cmd.Flag("searchText")
	assert.NotNil(t, searchFlag)
	assert.Equal(t, "", searchFlag.DefValue)

	pageSizeFlag := cmd.Flag("pageSize")
	assert.NotNil(t, pageSizeFlag)
	assert.Equal(t, "100", pageSizeFlag.DefValue)
}
