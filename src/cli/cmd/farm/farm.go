package farm

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/farm/dependecies"
	"github.com/cldcvr/terrarium/src/cli/cmd/farm/mappings"
	"github.com/cldcvr/terrarium/src/cli/cmd/farm/modules"
	"github.com/cldcvr/terrarium/src/cli/cmd/farm/resources"
	"github.com/spf13/cobra"
)

var farmCmd = &cobra.Command{
	Use:   "farm",
	Short: "Scrapes Terraform providers, modules, and mappings from the farm directory",
	Long:  `The 'farm' command provides subcommands to scrape Terraform providers, modules, and mappings from the specified farm directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display help or default action
		cmd.Help()
	},
}

func GetCmd() *cobra.Command {
	return farmCmd
}

func init() {
	farmCmd.AddCommand(resources.GetCmd())
	farmCmd.AddCommand(modules.GetCmd())
	farmCmd.AddCommand(mappings.GetCmd())
	farmCmd.AddCommand(dependecies.GetCmd())
}
