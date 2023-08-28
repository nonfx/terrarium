package resources

import (
	"path"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/cli"
	"github.com/cldcvr/terrarium/src/pkg/tf/runner"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	flagSchemaFile     string
	flagModuleListFile string
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
	cmd.Flags().StringVarP(&flagSchemaFile, "schema-file", "s", DefaultSchemaPath, "terraform provider schema json file path")
	cmd.Flags().StringVarP(&flagModuleListFile, "module-list-file", "f", "", "list file of modules to process")
	cmd.RunE = cmdRunE
}

func GetCmd() *cobra.Command {
	return cmd
}

func cmdRunE(cmd *cobra.Command, _ []string) error {
	// Connect to the database
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrapf(err, "error connecting to the database")
	}

	if flagModuleListFile == "" {
		cmd.Printf("Loading modules from the provider schema JSON file at '%s'...\n", flagSchemaFile)
		return loadFrom(g, flagSchemaFile)
	}

	cmd.Printf("Loading modules from modules list YAML file '%s'...\n", flagModuleListFile)
	moduleList, err := cli.LoadFarmModules(flagModuleListFile)
	if err != nil {
		return err
	}

	tfRunner := runner.NewTerraformRunner()
	for _, item := range moduleList.Farm {
		dir, _, err := item.CreateTerraformFile()
		if err != nil {
			return err
		}

		schemaFilePath := path.Join(dir, DefaultSchemaPath)
		if err := runner.TerraformProviderSchema(tfRunner, dir, schemaFilePath); err != nil {
			return err
		}

		if err := loadFrom(g, schemaFilePath); err != nil {
			return err
		}
	}

	return nil
}

func loadFrom(g db.DB, schemaFilePath string) error {
	cmd.Printf("Loading providers from '%s'\n", schemaFilePath)

	// Load providers schema from file
	providersSchema, err := loadProvidersSchema(schemaFilePath)
	if err != nil {
		return eris.Wrap(err, heredoc.Docf(`
			error loading providers schema file. make sure the schema file is created by following the instructions in the command help.
				terraform init && terraform providers schema -json > %s
		`, schemaFilePath))
	}

	providerCount, resCount, attrCount, err := pushProvidersSchemaToDB(providersSchema, g)
	if err != nil {
		return eris.Wrapf(err, "error writing data to db")
	}

	cmd.Printf("Successfully added %d Providers, %d Resources, and %d Attributes.\n", providerCount, resCount, attrCount)

	return nil
}
