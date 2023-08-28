package dependencies

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var depIfaceDirectoryFlag string

var cmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Harvests dependencies from the given directory",
	Long: heredoc.Docf(`
		The 'dependencies' command is used to harvest dependency information from YAML or YML files located
		in a specified directory. It parses these files to extract dependency details and stores them in the database
		for further reference.

		To use this command, provide the path to the directory containing the YAML or YML files using the '--dir' flag.
		The command will recursively process all valid YAML files within the directory, extracting information such as
		taxonomy, title, description, inputs, and outputs. The extracted data is then stored in the database.

		Example usage:
  			terrarium dependencies --dir /path/to/yaml/files

		Please ensure that the provided directory contains valid YAML or YML files with the appropriate structure to avoid any errors.
		`),
}

func init() {
	cmd.Flags().StringVarP(&depIfaceDirectoryFlag, "dir", "d", ".", "path to dependency directory")
	cmd.RunE = cmdRunE
}

func GetCmd() *cobra.Command {
	return cmd
}

func cmdRunE(cmd *cobra.Command, args []string) error {
	// Connect to the database
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrapf(err, "error connecting to the database")
	}
	err = processYAMLFiles(g, depIfaceDirectoryFlag)
	if err != nil {
		return err
	}
	cmd.Printf("Dependency interfaces successfully updated to the db\n")
	return nil
}
