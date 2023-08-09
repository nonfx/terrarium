package mappings

import (
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/cli/internal/constants"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	flagTFDir string
)

var cmd = &cobra.Command{
	Use:   "mappings",
	Short: "Scrapes resource attribute mappings from the terraform directory",
	Long:  "The 'mappings' command scrapes resource attribute mappings from the specified terraform directory.",
}

func init() {
	cmd.RunE = cmdRunE
	cmd.Flags().StringVarP(&flagTFDir, "dir", "d", ".", "terraform directory path")
}

func GetCmd() *cobra.Command {
	return cmd
}

func cmdRunE(cmd *cobra.Command, _ []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrapf(err, "error connecting to db")
	}

	resourceTypeByName = make(map[string]*db.TFResourceType)
	cmd.Printf("Loading modules from '%s'...\n", flagTFDir)
	configs, _, err := tfconfig.LoadModulesFromResolvedSchema(filepath.Join(flagTFDir, constants.ModuleSchemaFilePath))
	if err != nil {
		return eris.Wrapf(err, "error loading module")
	}
	moduleCount := len(configs)
	log.Infof("Loaded %d modules\n", moduleCount)

	totalResourceDeclarationsCount := 0
	totalMappingsCreatedCount := 0
	for _, config := range configs {
		mappings, resourceCount, err := createMappingsForModule(g, config)
		if err != nil {
			return eris.Wrapf(err, "error create mappings")
		}
		totalResourceDeclarationsCount += resourceCount
		totalMappingsCreatedCount += len(mappings)
	}

	cmd.Printf("Processed %d resource declarations in %d modules and created %d mappings...\n", totalResourceDeclarationsCount, moduleCount, totalMappingsCreatedCount)
	return nil
}