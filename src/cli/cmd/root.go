package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terrarium/src/cli/cmd/generate"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest"
	"github.com/cldcvr/terrarium/src/cli/cmd/platform"
	"github.com/cldcvr/terrarium/src/cli/cmd/version"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagCfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "terrarium [command]",
	Short: "Terrarium is a set of tools for cloud infrastructure provisioning",
	Long:  `Terrarium is a set of tools meant to simplify cloud infrastructure provisioning. It provides tools for both app developers and DevOps teams. Terrarium helps DevOps teams in writing Terraform code and helps app developer teams in declaring app dependencies to generate working Terraform code.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(harvest.GetCmd())
	rootCmd.AddCommand(platform.GetCmd())
	rootCmd.AddCommand(generate.GetCmd())
	rootCmd.AddCommand(version.GetCmd())
	rootCmd.PersistentFlags().StringVar(&flagCfgFile, "config", "", "config file (default is $HOME/.terrarium.yaml)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if flagCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagCfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".terrarium" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".terrarium")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file", "file", viper.ConfigFileUsed())
	}

	config.LoadDefaults()
}
