package platform

import (
	"github.com/cldcvr/terrarium/src/cli/cmd/platform/lint"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "platform",
	Aliases: []string{"p"},
	Short:   "Terrarium platform template commands",
	Long:    "Commands to manage Terrarium platform template definitions",
}

func init() {
	cmd.AddCommand(lint.GetCmd())
}

func GetCmd() *cobra.Command {
	return cmd
}
