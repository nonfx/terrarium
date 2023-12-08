// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package mappings

import (
	"fmt"
	"path/filepath"

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
	cmd                *cobra.Command
	flagTFDir          string
	flagModuleListFile string
	flagWorkDir        string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "mappings",
		Short: "Scrapes resource attribute mappings from the terraform directory",
		Long: heredoc.Doc(
			`The 'mappings' command scrapes resource attribute mappings from the specified terraform directory.
			It parses Terraform code and its modules to find mappings between input and output resource attributes,
			such as linking an input attribute of one resource to an output attribute of another.
		`),
		RunE: cmdRunE,
	}

	cmd.Flags().StringVarP(&flagTFDir, "dir", "d", ".", "Path to the Terraform directory")
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
	for _, item := range moduleList.Groups() {
		log.Info("harvesting mappings", "groupName", item.Name)
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
	fmt.Fprintf(cmd.OutOrStdout(), "Loading modules from '%s'...\n", dir)
	configs, _, err := tfconfig.LoadModulesFromResolvedSchema(filepath.Join(dir, constants.ModuleSchemaFilePath))
	if err != nil {
		return eris.Wrapf(err, "error loading module")
	}
	moduleCount := len(configs)
	log.Infof("Loaded %d modules", moduleCount)

	totalResourceDeclarationsCount := 0
	totalMappingsCreatedCount := 0
	for _, config := range configs {
		mappings, resourceCount, err := createMappingsForModule(g, config, resourceTypeByName)
		if err != nil {
			return eris.Wrapf(err, "error create mappings")
		}
		totalResourceDeclarationsCount += resourceCount
		totalMappingsCreatedCount += len(mappings)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Processed %d resource declarations in %d modules and created %d mappings...\n", totalResourceDeclarationsCount, moduleCount, totalMappingsCreatedCount)

	return nil
}
