// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"fmt"
	"path"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/log"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/cli"
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

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "resources",
		Aliases: []string{"res"},
		Short:   "Harvests Terraform resources and attributes using the provider schema json",
		Long: heredoc.Docf(`
			Harvests Terraform resources and attributes using the provider schema json.

			This command requires terraform provider schema already generated. To do that, run:
				terraform init && terraform providers schema -json > %s
		`, DefaultSchemaPath),
		RunE: cmdRunE,
	}

	cmd.Flags().StringVarP(&flagSchemaFile, "schema-file", "s", DefaultSchemaPath, "terraform provider schema json file path")
	cmd.Flags().StringVarP(&flagModuleListFile, "module-list-file", "f", "", "list file of modules to process")
	cmd.Flags().StringVarP(&flagWorkDir, "workdir", "w", "", "store all module sources in this directory; improves performance by reusing data between harvest commands")

	return cmd
}

func cmdRunE(cmd *cobra.Command, _ []string) error {
	// Connect to the database
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrapf(err, "error connecting to the database")
	}

	if flagModuleListFile == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from the provider schema JSON file at '%s'...\n", flagSchemaFile)
		return loadFrom(g, flagSchemaFile)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from modules list YAML file '%s'...\n", flagModuleListFile)
	moduleList, err := cli.LoadFarmModules(flagModuleListFile)
	if err != nil {
		return err
	}

	tfRunner := runner.NewTerraformRunner()
	for _, item := range moduleList.Farm {
		log.Info("harvesting resources from module", "name", item.Name, "source", item.Source)
		dir, _, err := item.CreateTerraformFile(flagWorkDir)
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
