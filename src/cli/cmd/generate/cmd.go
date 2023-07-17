package generate

import (
	"github.com/spf13/cobra"
)

var (
	platformDirPath string
	appPaths        []string
	outputDirPath   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates Terraform code",
	Long:  `The 'generate' command generates Terraform code based on a platform and app dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func GetCmd() *cobra.Command {
	return generateCmd
}

func init() {
	generateCmd.PersistentFlags().StringVar(&platformDirPath, "platform-dir", ".", "Path to the platform directory")
	generateCmd.PersistentFlags().StringArrayVarP(&appPaths, "app-path", "a", []string{}, "Path(s) to the app directory or terrarium.yaml file")
	generateCmd.PersistentFlags().StringVarP(&outputDirPath, "output-dir", "o", ".terrarium", "Path to the output directory (default: .terrarium)")
}
