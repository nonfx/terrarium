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
	Short: "Terrarium platform template commands",
	Long:  "Commands to manage Terrarium platform template definitions",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	platformCmd.AddCommand(lint.GetCmd())
}
