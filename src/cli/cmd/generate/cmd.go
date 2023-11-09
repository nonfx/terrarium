// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"fmt"
	"os"
	"path"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/cli/internal/constants"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/cldcvr/terrarium/src/pkg/metadata/utils"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	flagPlatformDir         string
	flagOutDir              string
	flagApps                []string
	flagProfile             string
	flagIgnoreUnimplemented bool
	flagSkipEnvFile         bool
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "generate",
		Short: "Terrarium generate terraform code using platform and app dependencies",
		Long:  "Terrarium generate command composes the working terraform code using the given platform template and a set of dependencies",
		RunE:  cmdRunE,
	}

	cmd.Flags().StringVarP(&flagPlatformDir, "platform-dir", "p", ".", "path to the directory containing the Terrarium platform template")
	cmd.Flags().StringArrayVarP(&flagApps, "app", "a", nil, "path to the app directory or the app yaml file. can be more then one")
	cmd.Flags().StringVarP(&flagOutDir, "output-dir", "o", "./.terrarium", "path to the directory where you want to generate the output")
	cmd.Flags().StringVarP(&flagProfile, "configuration-profile", "c", "", "name of platform configuration profile to apply")
	cmd.Flags().BoolVar(&flagIgnoreUnimplemented, "ignore-unimplemented", false, "set this to ignore errors when a component is not implemented in the platform") // not recommended
	cmd.Flags().BoolVar(&flagSkipEnvFile, "skip-env-file", false, "set this to skip creating the env files for each app")                                         // not recommended

	return cmd
}

func cmdRunE(cmd *cobra.Command, args []string) error {
	if len(flagApps) == 0 {
		return eris.New("No Apps provided. use -a flag to set apps")
	}

	apps, err := fetchApps(flagApps)
	if err != nil {
		return err
	}

	m, _ := tfconfig.LoadModule(flagPlatformDir, &tfconfig.ResolvedModulesSchema{})

	existingYaml, _ := os.ReadFile(path.Join(flagPlatformDir, defaultYAMLFileName))

	pm, _ := platform.NewPlatformMetadata(m, existingYaml)

	err = utils.MatchAppAndPlatform(pm, apps, flagIgnoreUnimplemented)
	if err != nil {
		return err
	}

	err = os.MkdirAll(flagOutDir, constants.ReadWriteExecutePermissions)
	if err != nil {
		return eris.Wrapf(err, "failed to create directory for %s", flagOutDir)
	}

	blockCount, err := writeTF(pm.Graph, flagOutDir, apps, m, flagProfile)
	if err != nil {
		return eris.Wrapf(err, "failed to write terraform code to dir: %s", flagOutDir)
	}

	if !flagSkipEnvFile {
		err = writeAppsEnv(pm, apps)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully pulled %d of %d terraform blocks at: %s\n", blockCount, len(pm.Graph), flagOutDir)
	return nil
}
