package harvest

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/mappings"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/modules"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/resources"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "harvest",
	Aliases: []string{"farm"}, // DEPRECATED. Just for backward compatibility.
	Short:   "Commands to scans the Terrarium farm code to update the core terrarium DB.",
	Long:    `The 'harvest' command provides subcommands to scans Terraform providers, modules, and mappings from the specified farm directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display help or default action
		cmd.Help()
	},
}

func GetCmd() *cobra.Command {
	return cmd
}

func init() {
	cmd.AddCommand(resources.GetCmd())
	cmd.AddCommand(modules.GetCmd())
	cmd.AddCommand(mappings.GetCmd())
}
