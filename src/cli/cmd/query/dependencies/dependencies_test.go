// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package dependencies

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
)

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

func Test_fetchDependencies(t *testing.T) {
	config.LoadDefaults()
	filterOptionType := "db.FilterOption"
	mockdir := "/Users/xyz/abc/tf-dir"
	clitest := clitesting.CLITest{
		CmdToTest: NewCmd,
	}
	mockUuid1 := uuid.New()
	mockDB := &mocks.DB{}
	mockDB.On("QueryDependencies", mock.AnythingOfType(filterOptionType), mock.AnythingOfType(filterOptionType)).
		Return(nil, fmt.Errorf("mock error")).Once()
	mockDB.On("QueryDependencies", mock.AnythingOfType(filterOptionType), mock.AnythingOfType(filterOptionType)).
		Return(db.DependencyOutputs{
			{
				Dependency: db.Dependency{
					Model:       db.Model{ID: mockUuid1},
					InterfaceID: "test_id",
					Title:       "test_title",
					Description: "testing",
				},
			},
		}, nil)
	config.SetDBMocks(mockDB)
	tests := []clitesting.CLITestCase{
		{
			Name:    "db failure",
			WantErr: true,
		},
		{
			Name: "list dependency interface in json format",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {

				args := []string{"-o", "json", "", mockdir}

				cmd.SetArgs(args)
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				expectedOutput := "{\"dependencies\":[{\"interfaceId\":\"test_id\",\"title\":\"test_title\",\"description\":\"testing\",\"inputs\":null,\"outputs\":null}],\"page\":{\"size\":100,\"index\":0,\"total\":0}}"
				return assert.JSONEq(t, expectedOutput, string(output))
			},
		},
		{
			Name: "list dependency interface in tabular format",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				args := []string{"-o", "table", "", mockdir}
				cmd.SetArgs(args)
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				expectedOutput := `+--------------+------------+-------------+--------+---------+
| INTERFACE ID |   TITLE    | DESCRIPTION | INPUTS | OUTPUTS |
+--------------+------------+-------------+--------+---------+
| test_id      | test_title | testing     | N/A    | N/A     |
+--------------+------------+-------------+--------+---------+
`

				return assert.Equal(t, expectedOutput, string(output))
			},
		},
		{
			Name: "list dependency interface with pagesize",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				args := []string{"-o", "json", "", mockdir, "--pageSize", "5"}
				cmd.SetArgs(args)
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				expected := map[string]interface{}{
					"dependencies": []interface{}{
						map[string]interface{}{
							"description": "testing",
							"interfaceId": "test_id",
							"inputs":      nil,
							"outputs":     nil,
							"title":       "test_title",
						},
					},
					"page": map[string]interface{}{
						"size":  5,
						"index": 0,
						"total": 0,
					},
				}

				expectedJSON, err := json.Marshal(expected)
				if err != nil {
					t.Fatalf("Failed to marshal expected output to JSON: %v", err)
				}
				return assert.JSONEq(t, string(expectedJSON), string(output))
			},
		},
	}
	clitest.RunTests(t, tests)
}

func TestSchemaToString(t *testing.T) {
	tests := []struct {
		name     string
		schema   *terrariumpb.JSONSchema
		expected string
	}{
		{
			name:     "schema is nil",
			schema:   nil,
			expected: "N/A",
		},
		{
			name:     "schema properties is nil",
			schema:   &terrariumpb.JSONSchema{Type: "object", Title: "Test"},
			expected: "Type: object, Title: Test",
		},
		{
			name: "schema with properties but no default",
			schema: &terrariumpb.JSONSchema{
				Type:  "object",
				Title: "Test",
				Properties: map[string]*terrariumpb.JSONSchema{
					"prop1": {Type: "string", Title: "Prop1", Description: "Description1"},
				},
			},
			expected: "Type: object, Title: Test, Properties: {prop1: {Type: string, Title: Prop1, Description: Description1}}",
		},
		{
			name: "schema with properties and default",
			schema: &terrariumpb.JSONSchema{
				Type:  "object",
				Title: "Test",
				Properties: map[string]*terrariumpb.JSONSchema{
					"prop1": {Type: "string", Title: "Prop1", Description: "Description1", Default: structpb.NewStringValue("default1")},
				},
			},
			expected: "Type: object, Title: Test, Properties: {prop1: {Type: string, Title: Prop1, Description: Description1, Default: default1}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schemaToString(tt.schema)
			assert.Equal(t, tt.expected, result)
		})
	}
}
