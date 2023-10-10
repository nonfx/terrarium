// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package taxonomy

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
	flags            *terrariumpb.ListTaxonomyRequest
	flagOutputFormat string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "taxonomy",
		Aliases: []string{"t"},
		Short:   "Query available taxonomy",
		Long:    "command to query available taxonomy.",
		RunE:    queryTaxonomy,
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
