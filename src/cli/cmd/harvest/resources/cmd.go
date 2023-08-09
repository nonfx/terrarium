package resources

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	flagSchemaFile string
)

var cmd = &cobra.Command{
	Use:     "resources",
	Aliases: []string{"res"},
	Short:   "Harvests Terraform resources and attributes using the provider schema json",
	Long: heredoc.Docf(`
		Harvests Terraform resources and attributes using the provider schema json.

		This command requires terraform provider schema already generated. To do that, run:
			terraform init && terraform providers schema -json > %s
	`, DefaultSchemaPath),
}

func init() {
	cmd.Flags().StringVarP(&flagSchemaFile, "file", "f", DefaultSchemaPath, "terraform provider schema json file path")
	cmd.RunE = cmdRunE
}

func GetCmd() *cobra.Command {
	return cmd
}

func cmdRunE(cmd *cobra.Command, _ []string) error {

	cmd.Printf("Loading providers from '%s'\n", flagSchemaFile)

	// Load providers schema from file
	providersSchema, err := loadProvidersSchema(flagSchemaFile)
	if err != nil {
		return eris.Wrap(err, heredoc.Docf(`
			error loading providers schema file. make sure the schema file is created by following the instructions in the command help.
				terraform init && terraform providers schema -json > %s
		`, flagSchemaFile))
	}

	// Connect to the database
	db, err := config.DBConnect()
	if err != nil {
		return eris.Wrapf(err, "error connecting to the database")
	}

	providerCount, resCount, attrCount, err := pushProvidersSchemaToDB(providersSchema, db)
	if err != nil {
		return eris.Wrapf(err, "error writing data to db")
	}

	cmd.Printf("Successfully added %d Providers, %d Resources, and %d Attributes.\n", providerCount, resCount, attrCount)

	return nil
}
