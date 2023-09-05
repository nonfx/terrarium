package dependencies

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	flagOutputFormat string
	depIfaceIDFlag   string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "dependencies",
		Short: "List dependency details matching the component pattern",
		Long:  "command to list matching dependencies as per the filter chosen. provides variety of filters to list desired dependencies",
		RunE:  fectchDependencies,
	}

	cmd.Flags().StringVarP(&depIfaceIDFlag, "id", "i", "", "id of the dependency interface")
	cmd.Flags().StringVarP(&flagOutputFormat, "output", "o", "json", "Output format (json or table)")

	return cmd
}

func fectchDependencies(cmd *cobra.Command, args []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return err
	}
	result, err := g.FetchDependencyByInterfaceID(depIfaceIDFlag)
	if err != nil {
		return err
	}

	if flagOutputFormat == "json" {
		b, err := json.MarshalIndent(result, "", " ")
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), string(b))
	} else {
		displayInTable(result)
	}
	return nil
}

func displayInTable(dependency *terrariumpb.Dependency) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Interface ID", "Title", "Description", "Inputs", "Outputs"})

	table.Append([]string{
		dependency.InterfaceId,
		dependency.Title,
		dependency.Description,
		string(dependency.Inputs),
		string(dependency.Outputs),
	})

	table.Render()
}
