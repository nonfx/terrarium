package generate

import (
	"path/filepath"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	flagPlatformDir string
	flagOutDir      string
	flagComponents  []string
)

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "Terrarium generate terraform code using platform and app dependencies",
	Long:  "Terrarium generate command composes the working terraform code using the given platform template and a set of dependencies",
}

func init() {
	cmd.RunE = cmdRunE
	cmd.Flags().StringVarP(&flagPlatformDir, "platform-dir", "p", ".", "path to the directory containing the Terrarium platform template")
	cmd.Flags().StringVarP(&flagOutDir, "output-dir", "o", "./.terrarium", "path to the directory where you want to generate the output")
	cmd.Flags().StringArrayVarP(&flagComponents, "component", "c", nil, "name of the components to pull from the platform. Use one or more times.")
}

func GetCmd() *cobra.Command {
	return cmd
}

func cmdRunE(cmd *cobra.Command, args []string) error {
	if len(flagComponents) == 0 {
		return eris.New("No components provided. use -c flag to set components")
	}
	m, diags := tfconfig.LoadModule(flagPlatformDir, &tfconfig.ResolvedModulesSchema{})
	if diags.HasErrors() {
		absPath, _ := filepath.Abs(flagPlatformDir) //os.Getwd()
		return eris.Wrapf(diags.Err(), "failed to parse the given platform terraform module at: '%s' (%s)", flagPlatformDir, absPath)
	}

	pm, _ := platform.NewPlatformMetadata(m, nil)

	bIds := blocksToPull(pm.Graph, flagComponents...)
	blockCount, err := writeTF(pm.Graph, flagOutDir, bIds, m)
	if err != nil {
		return eris.Wrapf(err, "failed to write terraform code to dir: %s", flagOutDir)
	}

	cmd.Printf("Successfully pulled %d of %d terraform blocks at: %s\n", blockCount, len(pm.Graph), flagOutDir)

	return nil
}
