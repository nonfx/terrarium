// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/cli/internal/utils"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	cmd              *cobra.Command
	flags            *terrariumpb.ListPlatformsRequest
	flagOutputFormat string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:     "platforms",
		Aliases: []string{"p"},
		Short:   "Query available platform templates",
		Long: heredoc.Doc(`
			The 'platforms' command allows you to query and list available platform templates in the system.
			It offers various filtering and pagination options, and the output can be formatted as either a table or JSON.

			- Taxonomy Flag (-t, --taxonomy):
			  The taxonomy flag allows you to query platforms based on the taxonomy levels of the dependency interfaces implemented in the platform components. If a platform doesn't have any components that belong to the specified taxonomy levels, it will be omitted from the response.

			- Dependency Interface Flag (-i, --interface):
			  The interface flag enables you to query platforms that implement certain dependency interfaces. Platforms that do not implement the given dependency interface UUIDs will be omitted. The UUIDs for dependency interfaces can be found using the 'query dependencies' command.

			Example Usage:
			  platforms --pageSize=50 --pageIndex=1 -t "storage/database" -o json
		`),
		RunE: queryPlatforms,
	}

	flags = &terrariumpb.ListPlatformsRequest{
		Page:          &terrariumpb.Page{},
		InterfaceUuid: []string{},
	}

	cmd.Flags().Int32Var(&flags.Page.Size, "pageSize", 100, "page size flag")
	cmd.Flags().Int32Var(&flags.Page.Index, "pageIndex", 0, "page index flag")
	cmd.Flags().StringVarP(&flags.Taxonomy, "taxonomy", "t", "", "taxonomy levels joined by `/`")
	cmd.Flags().StringVarP(&flags.Search, "search", "s", "", "search by platform name or repository")
	cmd.Flags().StringArrayVarP(&flags.InterfaceUuid, "interface", "i", []string{}, "search by dependency interface uuid implemented in the platform")
	cmd.Flags().StringVarP(&flagOutputFormat, "output", "o", "table", "Output format (json or table)")

	return cmd
}

func queryPlatforms(cmd *cobra.Command, args []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrap(err, "error connecting to the database")
	}

	result, err := g.QueryPlatforms(db.PlatformRequestToFilters(flags)...)
	if err != nil {
		return eris.Wrap(err, "error running database query")
	}

	f := utils.OutputFormatter[*terrariumpb.ListPlatformsResponse, *terrariumpb.Platform]{
		Writer: cmd.OutOrStdout(),
		Data: &terrariumpb.ListPlatformsResponse{
			Page:      flags.Page,
			Platforms: result.ToProto(),
		},
		RowHeaders: []string{"ID", "Title", "Repo", "Commit", "Components"},
		Array:      func(ltr *terrariumpb.ListPlatformsResponse) []*terrariumpb.Platform { return ltr.Platforms },
		Row: func(e *terrariumpb.Platform) []string {
			shortSha := e.RepoCommit
			if len(e.RepoCommit) > 7 {
				shortSha = e.RepoCommit[:7]
			}
			return []string{e.Id, e.Title, e.RepoUrl, shortSha, fmt.Sprintf("%d", e.Components)}
		},
	}

	return f.WriteJsonOrTable(flagOutputFormat == "json")
}
