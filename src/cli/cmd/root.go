package cmd

import (
	"fmt"
	"os"

	"github.com/cldcvr/terrarium/src/cli/cmd/farm"
	"github.com/cldcvr/terrarium/src/cli/cmd/generate"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "terrarium",
	Short: "Terrarium is a set of tools for cloud infrastructure provisioning",
	Long:  `Terrarium is a set of tools meant to simplify cloud infrastructure provisioning. It provides tools for both app developers and DevOps teams. Terrarium helps DevOps teams in writing Terraform code and helps app developer teams in declaring app dependencies to generate working Terraform code.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display help or default action
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(farm.GetCmd())
	rootCmd.AddCommand(generate.GetCmd())
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.terrarium.yaml)")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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
		fmt.Fprintf(rootCmd.OutOrStderr(), "Using config file: %s\n", viper.ConfigFileUsed())
	}

	config.LoadDefaults()
}
