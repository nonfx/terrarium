package harvest

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/dependencies"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/mappings"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/modules"
	"github.com/cldcvr/terrarium/src/cli/cmd/harvest/resources"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "harvest",
	Aliases: []string{"h"},
	Short:   "Commands to scans the Terrarium farm code and update the core terrarium DB.",
	Long:    `The 'harvest' command provides subcommands to scan Terraform providers, modules, mappings, and more from the specified farm directory.`,
}

func init() {
	cmd.AddCommand(resources.GetCmd())
	cmd.AddCommand(modules.GetCmd())
	cmd.AddCommand(mappings.GetCmd())
	cmd.AddCommand(dependencies.GetCmd())
}

func GetCmd() *cobra.Command {
	return cmd
}
