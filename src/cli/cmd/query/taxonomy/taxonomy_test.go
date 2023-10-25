// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package taxonomy

import (
	"context"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	dbmocks "github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

func TestNewCmd(t *testing.T) {
	testSetup := clitesting.CLITest{
		CmdToTest: NewCmd,
	}

	testSetup.RunTests(t, []clitesting.CLITestCase{
		{
			Name: "error connecting to db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				config.SetDBMocks(nil)
			},
			Args:     []string{},
			WantErr:  true,
			ExpError: "error connecting to the database: mocked err: connection failed",
		},
		{
			Name: "db query error",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("QueryTaxonomies", mock.Anything).Return(nil, eris.New("mock error"))
				config.SetDBMocks(mockedDB)
			},
			Args:     []string{},
			WantErr:  true,
			ExpError: "error running database query",
		},
		{
			Name: "success default",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("QueryTaxonomies", mock.Anything, mock.Anything).Return(db.Taxonomies{
					*db.TaxonomyFromLevels("mocked-l1", "l2", "l3"),
				}, nil)
				config.SetDBMocks(mockedDB)
			},
			Args:           []string{"-t", "mocked-l1/l2"},
			ValidateOutput: clitesting.ValidateOutputMatch("  #  ID                                    TAXONOMY         \n  1  00000000-0000-0000-0000-000000000000  mocked-l1/l2/l3  \n\nPage: 1 of 0 | Page Size: 100\n"),
		},
		{
			Name: "success json",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("QueryTaxonomies", mock.Anything).Return(db.Taxonomies{
					*db.TaxonomyFromLevels("mocked-l1", "l2", "l3"),
				}, nil)
				config.SetDBMocks(mockedDB)
			},
			Args:           []string{"-o", "json"},
			ValidateOutput: clitesting.ValidateOutputJson(`{"taxonomy":[{"id":"00000000-0000-0000-0000-000000000000","levels":["mocked-l1","l2","l3"]}],"page":{"size":100,"index":0,"total":0}}`),
		},
	})
}
