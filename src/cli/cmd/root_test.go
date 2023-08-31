// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/spf13/cobra"
)

func TestCmd(t *testing.T) {
	clitest := clitesting.CLITest{
		CmdToTest: func() *cobra.Command {
			cmd := newCmd()
			cmd.AddCommand(&cobra.Command{
				Use: "testmock",
				RunE: func(cmd *cobra.Command, args []string) error {
					cmd.Print(config.LogLevel())
					return nil
				},
			})

			return cmd
		},
	}

	clitest.RunTests(t, []clitesting.CLITestCase{
		{
			Name:           "render help",
			ValidateOutput: clitesting.ValidateOutputContains("Terrarium is a set of tools meant to simplify cloud infrastructure provisioning. It provides tools for both app developers and DevOps teams. Terrarium helps DevOps teams in writing Terraform code and helps app developer teams in declaring app dependencies to generate working Terraform code.\n\nUsage:\n  terrarium [command]"),
		},
		{
			Name: "call testmock command",
			Args: []string{"testmock"},
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				t.Setenv("TR_LOG_LEVEL", "mocked1")
			},
			ValidateOutput: clitesting.ValidateOutputMatch("mocked1"),
		},
		{
			Name: "invalid config path",
			Args: []string{"testmock", "--config", "./invalid-path"},
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				t.Setenv("TR_LOG_LEVEL", "mocked2")
			},
			ValidateOutput: clitesting.ValidateOutputMatch("mocked2"),
		},
	})
}
