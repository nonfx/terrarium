// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package components

import (
	"context"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	dbmocks "github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

func TestNewCmd(t *testing.T) {
	testSetup := clitesting.CLITest{
		CmdToTest: NewCmd,
	}

	mockedDependencyData := db.PlatformComponent{
		Model:        db.Model{ID: uuid.MustParse("e87eae47-66d1-4131-9e64-a20d6c7baa0d")},
		PlatformID:   uuid.MustParse("1d3edff9-8bdb-45f4-bafb-20b5fab97545"),
		DependencyID: uuid.MustParse("ea2bf723-363c-4414-b30f-f87a24c89592"),
		Dependency: db.Dependency{
			InterfaceID: "mock_dep",
			Title:       "mocked dependency",
			Description: "mocked description",
			Taxonomy:    &db.Taxonomy{Level1: "mocked-l1", Level2: "l2", Level3: "l3"},
			Attributes: db.DependencyAttributes{
				{
					Name:     "attr1",
					Computed: false,
				},
				{
					Name:     "attr2",
					Computed: true,
				},
			},
		},
	}

	testSetup.RunTests(t, []clitesting.CLITestCase{
		{
			Name: "error connecting to db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				config.SetDBMocks(nil)
			},
			Args:     []string{},
			WantErr:  true,
			ExpError: "invalid inputs",
		},
		{
			Name: "error connecting to db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				config.SetDBMocks(nil)
			},
			Args:     []string{"-p", "00000000-0000-0000-0000-000000000000"},
			WantErr:  true,
			ExpError: "error connecting to the database: mocked err: connection failed",
		},
		{
			Name: "db query error",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("QueryPlatformComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, eris.New("mock error"))
				config.SetDBMocks(mockedDB)
			},
			Args:     []string{"-p", "00000000-0000-0000-0000-000000000000"},
			WantErr:  true,
			ExpError: "error running database query",
		},
		{
			Name: "success default",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("QueryPlatformComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(db.PlatformComponents{mockedDependencyData}, nil)
				config.SetDBMocks(mockedDB)
			},
			Args:           []string{"-p", "00000000-0000-0000-0000-000000000000", "-t", "mocked-l1/l2"},
			ValidateOutput: clitesting.ValidateOutputMatch("  #  ID                                    DEPENDENCY UUID                       TITLE              DEPENDENCY  TAXONOMY         \n  1  e87eae47-66d1-4131-9e64-a20d6c7baa0d  ea2bf723-363c-4414-b30f-f87a24c89592  mocked dependency  mock_dep    mocked-l1/l2/l3  \n\nPage: 1 of 0 | Page Size: 100\n"),
		},
		{
			Name: "success json",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("QueryPlatformComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(db.PlatformComponents{mockedDependencyData}, nil)
				config.SetDBMocks(mockedDB)
			},
			Args:           []string{"-p", "00000000-0000-0000-0000-000000000000", "-o", "json"},
			ValidateOutput: clitesting.ValidateOutputJson(`{"components":[{"id":"e87eae47-66d1-4131-9e64-a20d6c7baa0d", "interfaceUuid":"ea2bf723-363c-4414-b30f-f87a24c89592", "interfaceId":"mock_dep", "taxonomy":["mocked-l1", "l2", "l3"], "title":"mocked dependency", "description":"mocked description", "inputs":{"title":"", "description":"", "type":"object", "default":null, "examples":[], "enum":[], "minLength":0, "maxLength":0, "pattern":"", "format":"", "minimum":0, "maximum":0, "exclusiveMinimum":false, "exclusiveMaximum":false, "multipleOf":0, "items":null, "additionalItems":false, "minItems":0, "maxItems":0, "uniqueItems":false, "properties":{"attr1":{"title":"", "description":"", "type":"", "default":null, "examples":[], "enum":[], "minLength":0, "maxLength":0, "pattern":"", "format":"", "minimum":0, "maximum":0, "exclusiveMinimum":false, "exclusiveMaximum":false, "multipleOf":0, "items":null, "additionalItems":false, "minItems":0, "maxItems":0, "uniqueItems":false, "properties":{}, "required":[]}}, "required":[]}, "outputs":{"title":"", "description":"", "type":"object", "default":null, "examples":[], "enum":[], "minLength":0, "maxLength":0, "pattern":"", "format":"", "minimum":0, "maximum":0, "exclusiveMinimum":false, "exclusiveMaximum":false, "multipleOf":0, "items":null, "additionalItems":false, "minItems":0, "maxItems":0, "uniqueItems":false, "properties":{"attr2":{"title":"", "description":"", "type":"", "default":null, "examples":[], "enum":[], "minLength":0, "maxLength":0, "pattern":"", "format":"", "minimum":0, "maximum":0, "exclusiveMinimum":false, "exclusiveMaximum":false, "multipleOf":0, "items":null, "additionalItems":false, "minItems":0, "maxItems":0, "uniqueItems":false, "properties":{}, "required":[]}}, "required":[]}}], "page":{"size":100, "index":0, "total":0}}`),
		},
	})
}
