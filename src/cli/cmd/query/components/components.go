// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/cli/internal/utils"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/taxonomy"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	cmd              *cobra.Command
	flags            *terrariumpb.ListComponentsRequest
	flagOutputFormat string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "components",
		Aliases: []string{"c"},
		Short:   "Query available components",
		Long:    "command to query available components.",
		RunE:    queryComponents,
	}

	flags = &terrariumpb.ListComponentsRequest{
		Page: &terrariumpb.Page{},
	}

	cmd.Flags().Int32Var(&flags.Page.Size, "pageSize", 100, "page size flag")
	cmd.Flags().Int32Var(&flags.Page.Index, "pageIndex", 0, "page index flag")
	cmd.Flags().StringVarP(&flags.Taxonomy, "taxonomy", "t", "", "taxonomy levels joined by `/`")
	cmd.Flags().StringVarP(&flags.Search, "search", "s", "", "search by component name or repository")
	cmd.Flags().StringVarP(&flags.PlatformId, "platformId", "p", "", "platform id to query in")
	cmd.Flags().StringVarP(&flagOutputFormat, "output", "o", "table", "Output format (json or table)")

	return cmd
}

func queryComponents(cmd *cobra.Command, args []string) error {
	err := flags.Validate()
	if err != nil {
		return eris.Wrap(err, "invalid inputs")
	}

	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrap(err, "error connecting to the database")
	}

	result, err := g.QueryPlatformComponents(db.ComponentRequestToFilters(flags)...)
	if err != nil {
		return eris.Wrap(err, "error running database query")
	}

	f := utils.OutputFormatter[*terrariumpb.ListComponentsResponse, *terrariumpb.Component]{
		Writer: cmd.OutOrStdout(),
		Data: &terrariumpb.ListComponentsResponse{
			Page:       flags.Page,
			Components: result.ToProto(),
		},
		RowHeaders: []string{"ID", "Dependency UUID", "Title", "Dependency", "Taxonomy"},
		Array:      func(ltr *terrariumpb.ListComponentsResponse) []*terrariumpb.Component { return ltr.Components },
		Row: func(e *terrariumpb.Component) []string {
			return []string{e.Id, e.InterfaceUuid, e.Title, e.InterfaceId, taxonomy.NewTaxonomy(e.Taxonomy...).String()}
		},
	}

	return f.WriteJsonOrTable(flagOutputFormat == "json")
}
