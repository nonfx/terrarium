package dependecies

import (
	"log"

	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {
	return dependencyCmd
}

var dependencyCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Harvests dependencies from the given directory",
	Long:  "Harvests dependencies from the directory and adds it to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func main() {
	log.Println("fetching dependencies...")
}
