package dependencies

import (
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/spf13/cobra"
)

var dependencyCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "List dependency details matching the component pattern",
	Long:  "command to list mathcing dependencies as per the filter chosen. provides variety of filters to list desired dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fectchDependencies(cmd, args)
	},
}

func GetCmd() *cobra.Command {
	return dependencyCmd
}

func fectchDependencies(cmd *cobra.Command, args []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return err
	}
	g.
	return nil
}
