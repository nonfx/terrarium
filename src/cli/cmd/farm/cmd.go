// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package farm

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/farm/update"
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "farm",
		Short: "terrarium farm related commands",
		Long:  `commands to access farm repository`,
	}

	cmd.AddCommand(update.NewCmd())

	return cmd
}
