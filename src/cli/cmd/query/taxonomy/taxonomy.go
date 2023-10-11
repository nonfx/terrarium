// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package taxonomy

import (
	"github.com/MakeNowJust/heredoc/v2"
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
	flags            *terrariumpb.ListTaxonomyRequest
	flagOutputFormat string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "taxonomy",
		Aliases: []string{"t"},
		Short:   "Query available taxonomy",
		Long: heredoc.Doc(`The 'taxonomy' command allows you to query and list available taxonomy values in the system.
		It provides various options to filter and paginate the results. You can specify the page size and index,
		as well as the taxonomy levels you are interested in. The output can be formatted as either a table or JSON.

		Example Usage:
			taxonomy --pageSize=50 --pageIndex=1 -t "storage/database" -o json
		`),
		RunE: queryTaxonomy,
	}

	flags = &terrariumpb.ListTaxonomyRequest{
		Page: &terrariumpb.Page{},
	}

	cmd.Flags().Int32Var(&flags.Page.Size, "pageSize", 100, "page size flag")
	cmd.Flags().Int32Var(&flags.Page.Index, "pageIndex", 0, "page index flag")
	cmd.Flags().StringVarP(&flags.Taxonomy, "taxonomy", "t", "", "taxonomy levels joined by `/`")
	cmd.Flags().StringVarP(&flagOutputFormat, "output", "o", "table", "Output format (json or table)")

	return cmd
}

func queryTaxonomy(cmd *cobra.Command, args []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrap(err, "error connecting to the database")
	}

	result, err := g.QueryTaxonomies(db.TaxonomyRequestToFilters(flags)...)
	if err != nil {
		return eris.Wrap(err, "error running database query")
	}

	f := utils.OutputFormatter[*terrariumpb.ListTaxonomyResponse, *terrariumpb.Taxonomy]{
		Writer: cmd.OutOrStdout(),
		Data: &terrariumpb.ListTaxonomyResponse{
			Page:     flags.Page,
			Taxonomy: result.ToProto(),
		},
		RowHeaders: []string{"ID", "Taxonomy"},
		Array:      func(ltr *terrariumpb.ListTaxonomyResponse) []*terrariumpb.Taxonomy { return ltr.Taxonomy },
		Row: func(t *terrariumpb.Taxonomy) []string {
			return []string{t.Id, taxonomy.NewTaxonomy(t.Levels...).String()}
		},
	}

	return f.WriteJsonOrTable(flagOutputFormat == "json")
}
