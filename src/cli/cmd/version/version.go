package version

import (
	"fmt"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/build"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Fprintf(cmd.OutOrStdout(), format(build.Version, build.Date))
		return
	},
}

func format(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	if buildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, buildDate)
	}

	return fmt.Sprintf("terrarium version %s\n", version)
}

func GetCmd() *cobra.Command {
	return versionCmd
}
