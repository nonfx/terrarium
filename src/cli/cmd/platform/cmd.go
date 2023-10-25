// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/platform/lint"
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "platform",
		Aliases: []string{"p"},
		Short:   "Terrarium platform template commands",
		Long:    "Commands to manage Terrarium platform template definitions",
	}

	cmd.AddCommand(lint.NewCmd())

	return cmd
}
