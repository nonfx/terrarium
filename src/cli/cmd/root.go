// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terrarium/src/cli/cmd/farm"
	"github.com/cldcvr/terrarium/src/cli/cmd/generate"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest"
	"github.com/cldcvr/terrarium/src/cli/cmd/platform"
	"github.com/cldcvr/terrarium/src/cli/cmd/query"
	"github.com/cldcvr/terrarium/src/cli/cmd/version"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd *cobra.Command

	flagCfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	newCmd()
}

func newCmd() *cobra.Command {
	rootCmd = &cobra.Command{
		Use:   "terrarium [command]",
		Short: "Terrarium is a set of tools for cloud infrastructure provisioning",
		Long:  `Terrarium is a set of tools meant to simplify cloud infrastructure provisioning. It provides tools for both app developers and DevOps teams. Terrarium helps DevOps teams in writing Terraform code and helps app developer teams in declaring app dependencies to generate working Terraform code.`,
	}

	rootCmd.AddCommand(harvest.NewCmd())
	rootCmd.AddCommand(platform.NewCmd())
	rootCmd.AddCommand(generate.NewCmd())
	rootCmd.AddCommand(version.NewCmd())
	rootCmd.AddCommand(query.NewCmd())
	rootCmd.AddCommand(farm.NewCmd())

	rootCmd.PersistentFlags().StringVar(&flagCfgFile, "config", "", "config file (default is $HOME/.terrarium/config.yaml)")

	return rootCmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Debugf("%+v", err)
		os.Exit(1)
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if flagCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagCfgFile)
	} else {
		const defaultConfigDirPath = "~/.terrarium/"

		// Resolve the path and create directory if not present.
		dirPath, err := utils.SetupDir(defaultConfigDirPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in `$HOME/.terrarium` directory with name "config" (without extension)
		viper.AddConfigPath(dirPath)
		viper.SetConfigName("config")
	}

	config.LoadDefaults()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		config.LoggerConfigDefault()
		log.Info("Using config file", "file", viper.ConfigFileUsed())
	}
}
