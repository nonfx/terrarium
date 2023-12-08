// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package modules

import (
	"fmt"
	"path"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/log"
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/cli/internal/constants"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/modulelist"
	"github.com/cldcvr/terrarium/src/pkg/tf/runner"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	flagTFDir          string
	flagIncludeLocal   bool
	flagModuleListFile string
	flagWorkDir        string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "modules",
		Aliases: []string{"mo"},
		Short:   "Harvests Terraform modules and attributes into the database",
		Long: heredoc.Doc(`
			The 'modules' command harvests all Terraform modules and their attributes into the database.

			This command can operate in two modes:
			1. Direct scraping from a specified Terraform directory.
			2. Processing a list of modules specified in a module list file.

			For direct scraping, ensure to run "terraform init" in the Terraform directory before using this command.
			When using a module list file, the command processes only the specified modules.

			Additional flags allow including local modules and specifying a working directory for improved performance.
		`),
		RunE: cmdRunE,
	}

	cmd.Flags().StringVarP(&flagTFDir, "dir", "d", ".", "Path to the Terraform directory")
	cmd.Flags().BoolVarP(&flagIncludeLocal, "enable-local-modules", "l", false, "Include local modules in the scraping process")
	cmd.Flags().StringVarP(&flagModuleListFile, "module-list-file", "f", "", "Path to a file listing modules to process. In this mode, 'terraform init' and 'terraform providers schema -json' are executed automatically. More details at https://github.com/cldcvr/terrarium/blob/main/src/pkg/metadata/modulelist/readme.md")
	cmd.Flags().StringVarP(&flagWorkDir, "workdir", "w", "", "Directory for storing module sources. Using a workdir improves performance by reusing data between harvesting multiple modules. This flag should be used in conjunction with 'module-list-file'.")

	return cmd
}

func cmdRunE(cmd *cobra.Command, _ []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrapf(err, "error connecting to the database")
	}

	if flagModuleListFile == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from provided directory '%s'...\n", flagTFDir)
		return loadFrom(g, flagTFDir)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from provided list file '%s'...\n", flagModuleListFile)
	moduleList, err := modulelist.LoadFarmModules(flagModuleListFile)
	if err != nil {
		return err
	}

	tfRunner := runner.NewTerraformRunner()
	for _, item := range moduleList.Groups().FilterExport(true) {
		log.Info("harvesting module", "groupName", item.Name)
		dir, err := item.CreateTerraformFile(flagWorkDir)
		if err != nil {
			return err
		}

		if err := runner.TerraformInit(tfRunner, dir); err != nil {
			return err
		}

		if err := loadFrom(g, dir); err != nil {
			return err
		}
	}

	return nil
}

func loadFrom(g db.DB, dir string) error {
	resourceTypeByName := make(map[string]*db.TFResourceType)

	// load modules
	fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from '%s'...\n", dir)

	filters := []tfconfig.ResolvedModuleSchemaFilter{tfconfig.FilterModulesOmitHidden}
	if !flagIncludeLocal {
		filters = append(filters, tfconfig.FilterModulesOmitLocal)
	}

	schemaFilePath := path.Clean(path.Join(dir, constants.ModuleSchemaFilePath))
	configs, _, err := tfconfig.LoadModulesFromResolvedSchema(schemaFilePath, filters...)
	if err != nil {
		return eris.Wrapf(err, "error loading module")
	}

	log.Info("Loaded modules", "count", len(configs))

	for _, config := range configs {
		log.Info("Processing module", "path", config.Path)

		moduleDB := toTFModule(config)
		// terraform generated schema file has an empty value.
		// check to avoid persisting empty values
		if len(strings.TrimSpace(moduleDB.ModuleName)) == 0 || len(strings.TrimSpace(moduleDB.Source)) == 0 {
			continue
		}
		if _, err := g.CreateTFModule(moduleDB); err != nil {
			return eris.Wrapf(err, "error creating module record")
		}

		for varName, v := range config.Variables {
			if varAttrReferences, ok := config.Inputs[varName]; ok { // found a resolution for this variable to resource attribute
				for varAttributePath, resourceReferences := range varAttrReferences {
					for _, res := range resourceReferences {
						if attr, err := createAttributeRecord(g, moduleDB, v, varAttributePath, res, resourceTypeByName); err != nil {
							return eris.Wrapf(err, "error creating module input-attribute")
						} else if attr == nil {
							continue
						}
					}
				}
			}
		}

		for _, o := range config.Outputs {
			if attr, err := createAttributeRecord(g, moduleDB, o, "", o.Value, resourceTypeByName); err != nil {
				return eris.Wrapf(err, "error creating module output-attribute")
			} else if attr == nil {
				continue
			}
		}
		log.Info("Done processing module", "path", config.Path)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Processed %d modules\n", len(configs))

	return nil
}
