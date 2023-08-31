// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/build"
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE:  cmdRunE,
	}

	return cmd
}

func cmdRunE(cmd *cobra.Command, args []string) (err error) {
	cmd.Print(format(build.Version, build.Date))
	return
}

func format(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	if buildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, buildDate)
	}

	return fmt.Sprintf("terrarium version %s\n", version)
}
