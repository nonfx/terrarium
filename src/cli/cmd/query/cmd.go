// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package query

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/platforms"
	"github.com/cldcvr/terrarium/src/cli/cmd/query/dependencies"
	"github.com/cldcvr/terrarium/src/cli/cmd/query/modules"
	"github.com/cldcvr/terrarium/src/cli/cmd/query/taxonomy"
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "List modules matching the source pattern",
		Long:    `commands to list matching modules as per the filter chosen. provides variety of filters to list desired modules`,
	}

	cmd.AddCommand(modules.NewCmd())
	cmd.AddCommand(dependencies.NewCmd())
	cmd.AddCommand(taxonomy.NewCmd())
	cmd.AddCommand(platforms.NewCmd())

	return cmd
}
