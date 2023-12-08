// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"fmt"
	"path"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/log"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/modulelist"
	"github.com/cldcvr/terrarium/src/pkg/tf/runner"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	flagSchemaFile     string
	flagModuleListFile string
	flagWorkDir        string
)

const DefaultSchemaPath = ".terraform/providers/schema.json"

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "resources",
		Aliases: []string{"res"},
		Short:   "Harvests Terraform providers, resource types, and resource attributes",
		Long: heredoc.Docf(`
			Harvests Terraform providers, resource types, and resource attributes.

			This command operates in two modes:
			1. Using a pre-generated provider schema JSON file.
			2. Using a module list file, where 'terraform init' and 'terraform providers schema -json'
			   are executed automatically for multiple given modules.

			For the first mode, ensure the provider schema JSON file is generated using:
				terraform init && terraform providers schema -json > %s

			In the second mode, only specify a module list file, and the necessary Terraform commands
			are run internally to generate the required data.
		`, DefaultSchemaPath),
		RunE: cmdRunE,
	}

	cmd.Flags().StringVarP(&flagSchemaFile, "schema-file", "s", DefaultSchemaPath, "Path to the Terraform provider schema JSON file. Use this for the first mode of operation.")
	cmd.Flags().StringVarP(&flagModuleListFile, "module-list-file", "f", "", "Path to a file listing modules to process. In this mode, 'terraform init' and 'terraform providers schema -json' are executed automatically. More details at https://github.com/cldcvr/terrarium/blob/main/src/pkg/metadata/modulelist/readme.md")
	cmd.Flags().StringVarP(&flagWorkDir, "workdir", "w", "", "Directory for storing module sources. Using a workdir improves performance by reusing data between harvesting multiple modules. This flag should be used in conjunction with 'module-list-file'.")

	return cmd
}

func cmdRunE(cmd *cobra.Command, _ []string) error {
	// Connect to the database
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrapf(err, "error connecting to the database")
	}

	// First mode - using schema file
	if flagModuleListFile == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from the provider schema JSON file at '%s'...\n", flagSchemaFile)
		return loadFrom(g, flagSchemaFile)
	}

	// Second mode using module list file

	fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from modules list YAML file '%s'...\n", flagModuleListFile)
	moduleList, err := modulelist.LoadFarmModules(flagModuleListFile)
	if err != nil {
		return err
	}

	tfRunner := runner.NewTerraformRunner()
	for _, item := range moduleList.Groups() {
		log.Info("harvesting resources", "groupName", item.Name)
		dir, err := item.CreateTerraformFile(flagWorkDir)
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
	fmt.Fprintf(cmd.OutOrStdout(), "Loading providers from '%s'\n", schemaFilePath)

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

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully added %d Providers, %d Resources, and %d Attributes.\n", providerCount, resCount, attrCount)

	return nil
}
