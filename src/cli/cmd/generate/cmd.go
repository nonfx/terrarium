package generate

import (
	"os"
	"path"
	"path/filepath"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	flagPlatformDir string
	flagOutDir      string
	flagApps        []string
)

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "Terrarium generate terraform code using platform and app dependencies",
	Long:  "Terrarium generate command composes the working terraform code using the given platform template and a set of dependencies",
}

func init() {
	cmd.RunE = cmdRunE
	cmd.Flags().StringVarP(&flagPlatformDir, "platform-dir", "p", ".", "path to the directory containing the Terrarium platform template")
	cmd.Flags().StringArrayVarP(&flagApps, "app", "a", nil, "path to the app directory or the app yaml file. can be more then one")
	cmd.Flags().StringVarP(&flagOutDir, "output-dir", "o", "./.terrarium", "path to the directory where you want to generate the output")
}

func GetCmd() *cobra.Command {
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

	m, diags := tfconfig.LoadModule(flagPlatformDir, &tfconfig.ResolvedModulesSchema{})
	if diags.HasErrors() {
		absPath, _ := filepath.Abs(flagPlatformDir)
		return eris.Wrapf(diags.Err(), "failed to parse the given platform terraform module at: '%s' (%s)", flagPlatformDir, absPath)
	}

	existingYaml, _ := os.ReadFile(path.Join(flagPlatformDir, defaultYAMLFileName))

	pm, _ := platform.NewPlatformMetadata(m, existingYaml)

	err = matchAppAndPlatform(pm, apps)
	if err != nil {
		return err
	}

	blockCount, err := writeTF(pm.Graph, flagOutDir, apps, m)
	if err != nil {
		return eris.Wrapf(err, "failed to write terraform code to dir: %s", flagOutDir)
	}

	cmd.Printf("Successfully pulled %d of %d terraform blocks at: %s\n", blockCount, len(pm.Graph), flagOutDir)

	return nil
}
