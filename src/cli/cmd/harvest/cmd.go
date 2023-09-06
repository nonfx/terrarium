// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package harvest

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/dependencies"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/mappings"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/modules"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/resources"
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "harvest",
		Aliases: []string{"h"},
		Short:   "Commands to scans the Terrarium farm code and update the core terrarium DB.",
		Long:    `The 'harvest' command provides subcommands to scan Terraform providers, modules, mappings, and more from the specified farm directory.`,
	}

	cmd.AddCommand(resources.NewCmd())
	cmd.AddCommand(modules.NewCmd())
	cmd.AddCommand(mappings.NewCmd())
	cmd.AddCommand(dependencies.NewCmd())

	return cmd
}
