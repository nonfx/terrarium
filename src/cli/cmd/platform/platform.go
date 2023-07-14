package platform

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/platform/lint"
	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {
	return platformCmd
}

var platformCmd = &cobra.Command{
	Use:   "platform",
	Short: "Defines reusable infrastructure components",
	Long:  "Commands to manage platform definitions.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	platformCmd.AddCommand(lint.GetCmd())
}
