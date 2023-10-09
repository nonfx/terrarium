// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package lint

import (
	"fmt"
	"os"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	flagDir string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "lint",
		Short: "Check a given directory is valid platform definition",
		Long:  "Analyze the directory and verify it constitutes a valid platform definition.",
		RunE:  cmdRunE,
	}

	cmd.Flags().StringVarP(&flagDir, "dir", "d", ".", "Path to platform directory to validate.")

	return cmd
}

func cmdRunE(cmd *cobra.Command, args []string) error {
	if err := checkDirExists(flagDir); err != nil {
		return err
	}
	err := lintPlatform(flagDir)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Platform parse and lint completed\n")
	return nil
}

func checkDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return eris.Wrapf(err, "could not open given directory '%s'", dir)
	}
	return nil
}
