// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package dependencies

import (
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/transporthelper"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	flagOutputFormat string
	flagPageSize     int32
	flagPageIndex    int32
	flagSearchText   string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "dependencies",
		Short: "List dependency details matching the dependency interface id",
		Long: heredoc.Docf(`
		The 'dependencies' command provides detailed information about the available dependencies.
		When invoked without any filters, it lists all the dependencies present. You can also refine
		the list using optional search criteria.

		Usage examples:
	  		terrarium dependencies                             // Lists all dependencies
	  		terrarium dependencies -s <SEARCH_TEXT>  // Filters the list based on the search text

		The above commands will fetch and display the details of the dependencies matching the criteria,
		if provided.
		`),
		RunE: fetchDependencies,
	}

	cmd.Flags().StringVarP(&flagSearchText, "searchText", "s", "", "optional search text")
	cmd.Flags().Int32Var(&flagPageSize, "pageSize", 100, "page size flag")
	cmd.Flags().Int32Var(&flagPageIndex, "pageIndex", 0, "page index flag")
	cmd.Flags().StringVarP(&flagOutputFormat, "output", "o", "json", "Output format (json or table)")

	return cmd
}

func fetchDependencies(cmd *cobra.Command, args []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return err
	}

	page := &terrariumpb.Page{
		Size:  flagPageSize,
		Index: flagPageIndex,
	}

	dbDependencies, err := g.QueryDependencies(
		db.DependencySearchFilter(flagSearchText),
		db.PaginateGlobalFilter(page.Size, page.Index, &page.Total),
	)
	if err != nil {
		return err
	}

	var result []*terrariumpb.Dependency

	// Populate the protobuf response with the result and page information
	pbRes := &terrariumpb.ListDependenciesResponse{
		Dependencies: dbDependencies.ToProto(),
		Page:         page,
	}

	if flagOutputFormat == "json" {
		b, err := transporthelper.CreateJSONBodyMarshaler().Marshaler.Marshal(pbRes)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), string(b))
	} else {
		displayInTable(cmd.OutOrStdout(), result)
	}
	return nil
}

func displayInTable(w io.Writer, dependencies []*terrariumpb.Dependency) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Interface ID", "Title", "Description", "Inputs", "Outputs"})

	for _, dependency := range dependencies {
		table.Append([]string{
			dependency.InterfaceId,
			dependency.Title,
			dependency.Description,
			schemaToString(dependency.Inputs),
			schemaToString(dependency.Outputs),
		})
	}

	table.Render()
}

// Utility function to convert JSONSchema into a readable string.
// This is a basic version; you can enhance it based on how detailed you want your output.
func schemaToString(schema *terrariumpb.JSONSchema) string {
	return fmt.Sprintf("Type: %s, Title: %s", schema.Type, schema.Title)
}
