// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package query

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/platforms"
	"github.com/cldcvr/terrarium/src/cli/cmd/query/dependencies"
	"github.com/cldcvr/terrarium/src/cli/cmd/query/modules"
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "query",
		Short: "List modules matching the source pattern",
		Long:  `commands to list matching modules as per the filter chosen. provides variety of filters to list desired modules`,
	}

	cmd.AddCommand(modules.NewCmd())
	cmd.AddCommand(dependencies.NewCmd())
	cmd.AddCommand(platforms.NewCmd())

	return cmd
}
