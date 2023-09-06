// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package dependencies

import (
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/transporthelper"
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
		Short: "List dependency details matching the dependency interface id",
		Long: heredoc.Docf(`
		The 'dependencies' command provides detailed information about dependencies based on their interface IDs.
		When invoked, it searches through the available dependencies and retrieves specific details related to the
		provided interface ID.

		Usage examples:
		  terrarium dependencies -i <specific_id>

		The above commands will fetch and list the details of the dependency associated with the specified interface ID.
		`),
		RunE: fectchDependencies,
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
		b, err := transporthelper.CreateJSONBodyMarshaler().Marshaler.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), string(b))
	} else {
		displayInTable(cmd.OutOrStdout(), result)
	}
	return nil
}

func displayInTable(w io.Writer, dependency *terrariumpb.Dependency) {
	table := tablewriter.NewWriter(w)
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
