package query

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/query/modules"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "query",
	Short: "List modules matching the source pattern",
	Long:  `commands to list mathcing modules as per the filter chosen. provides variety of filters to list desired modules`,
}

func GetCmd() *cobra.Command {
	return cmd
}

func init() {
	cmd.AddCommand(modules.GetCmd())
}
