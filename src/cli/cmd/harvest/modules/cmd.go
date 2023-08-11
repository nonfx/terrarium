package modules

import (
	"path"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/log"
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/cli/internal/constants"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/cli"
	"github.com/cldcvr/terrarium/src/pkg/tf/runner"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	flagTFDir          string
	flagIncludeLocal   bool
	flagModuleListFile string
)

var cmd = &cobra.Command{
	Use:     "modules",
	Aliases: []string{"mo"},
	Short:   "Scrapes Terraform modules and attributes from the terraform directory",
	Long: heredoc.Doc(`
		The 'modules' command scrapes all Terraform modules and its attributes from the specified terraform directory.

		Prerequisite: Run "terraform init" in the directory before using this command.
	`),
}

func init() {
	cmd.RunE = cmdRunE
	cmd.Flags().StringVarP(&flagTFDir, "dir", "d", ".", "terraform directory path")
	cmd.Flags().BoolVarP(&flagIncludeLocal, "enable-local-modules", "l", false, "A boolean flag to control include/exclude of local modules")
	cmd.Flags().StringVarP(&flagModuleListFile, "module-list-file", "f", "", "list file of modules to process")
}

func GetCmd() *cobra.Command {
	return cmd
}

func cmdRunE(cmd *cobra.Command, _ []string) error {
	g, err := config.DBConnect()
	if err != nil {
		panic(err)
	}

	if flagModuleListFile == "" {
		cmd.Printf("Loading modules from provided directory '%s'...\n", flagTFDir)
		return loadFrom(g, flagTFDir)
	}

	cmd.Printf("Loading modules from provided list file '%s'...\n", flagModuleListFile)
	moduleList, err := cli.LoadFarmModules(flagModuleListFile)
	if err != nil {
		return err
	}

	tfRunner := runner.NewTerraformRunner()
	for _, item := range moduleList.Farm {
		if item.Export {
			dir, _, err := item.CreateTerraformFile()
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
	}

	return nil
}

func loadFrom(g db.DB, dir string) error {
	resourceTypeByName := make(map[string]*db.TFResourceType)

	// load modules
	cmd.Printf("Loading modules from '%s'...\n", dir)

	filters := []tfconfig.ResolvedModuleSchemaFilter{tfconfig.FilterModulesOmitHidden}
	if !flagIncludeLocal {
		filters = append(filters, tfconfig.FilterModulesOmitLocal)
	}

	schemaFilePath := path.Clean(path.Join(dir, constants.ModuleSchemaFilePath))
	configs, _, err := tfconfig.LoadModulesFromResolvedSchema(schemaFilePath, filters...)
	if err != nil {
		panic(err)
	}

	log.Info("Loaded modules", "count", len(configs))

	for _, config := range configs {
		log.Info("Processing module", "path", config.Path)

		moduleDB := toTFModule(config)
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

	cmd.Printf("Processed %d modules\n", len(configs))

	return nil
}
