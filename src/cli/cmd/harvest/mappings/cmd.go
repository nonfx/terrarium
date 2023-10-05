// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package mappings

import (
	"fmt"
	"path/filepath"

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
	cmd                *cobra.Command
	flagTFDir          string
	flagModuleListFile string
	flagWorkDir        string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "mappings",
		Short: "Scrapes resource attribute mappings from the terraform directory",
		Long:  "The 'mappings' command scrapes resource attribute mappings from the specified terraform directory.",
		RunE:  cmdRunE,
	}

	cmd.Flags().StringVarP(&flagTFDir, "dir", "d", ".", "terraform directory path")
	cmd.Flags().StringVarP(&flagModuleListFile, "module-list-file", "f", "", "list file of modules to process")
	cmd.Flags().StringVarP(&flagWorkDir, "workdir", "w", "", "store all module sources in this directory; improves performance by reusing data between harvest commands")

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
	moduleList, err := cli.LoadFarmModules(flagModuleListFile)
	if err != nil {
		return err
	}

	tfRunner := runner.NewTerraformRunner()
	for _, item := range moduleList.Farm {
		dir, _, err := item.CreateTerraformFile(flagWorkDir)
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
