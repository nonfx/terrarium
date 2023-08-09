package lint

import "github.com/spf13/cobra"

var (
	flagDir string
)

var cmd = &cobra.Command{
	Use:   "lint",
	Short: "Check a given directory is valid platform definition",
	Long:  "Analyze the directory and verify it constitutes a valid platform definition.",
}

func init() {
	cmd.RunE = cmdRunE
	cmd.Flags().StringVar(&flagDir, "dir", ".", "Path to platform directory to validate.")
}

func GetCmd() *cobra.Command {
	return cmd
}

func cmdRunE(cmd *cobra.Command, args []string) error {
	err := lintPlatform(flagDir)
	if err != nil {
		return err
	}

	cmd.Printf("Platform parse and lint completed\n")
	return nil
}
