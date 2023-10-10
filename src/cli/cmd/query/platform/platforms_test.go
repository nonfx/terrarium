// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package platform

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
				mockedDB.On("QueryPlatforms", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, eris.New("mock error"))
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
				mockedDB.On("QueryPlatforms", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(db.Platforms{
					{
						Name:      "mock-pf1",
						CommitSHA: "0943637a8bf75489ef3e222db54c83a662c5f66a",
					},
				}, nil)
				config.SetDBMocks(mockedDB)
			},
			Args:           []string{"-t", "mocked-l1/l2"},
			ValidateOutput: clitesting.ValidateOutputMatch("  #  ID                                    TITLE     REPO  COMMIT   \n  1  00000000-0000-0000-0000-000000000000  mock-pf1        0943637  \n\nPage: 1 of 0 | Page Size: 100\n"),
		},
		{
			Name: "success json",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("QueryPlatforms", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(db.Platforms{
					{
						Name: "mock-pf1",
					},
				}, nil)
				config.SetDBMocks(mockedDB)
			},
			Args:           []string{"-o", "json"},
			ValidateOutput: clitesting.ValidateOutputJson(`{"platforms":[{"id":"00000000-0000-0000-0000-000000000000", "title":"mock-pf1", "description":"", "repoUrl":"", "repoDir":"", "repoCommit":"", "refLabel":"", "labelType":"gitRef_undefined", "components":0}], "page":{"size":100, "index":0, "total":0}}`),
		},
	})
}
