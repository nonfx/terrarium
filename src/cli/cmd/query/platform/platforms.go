// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platform

import (
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
		Short:   "Query available platforms",
		Long:    "command to query available platforms.",
		RunE:    queryPlatforms,
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
		RowHeaders: []string{"ID", "Title", "Repo", "Commit"},
		Array:      func(ltr *terrariumpb.ListPlatformsResponse) []*terrariumpb.Platform { return ltr.Platforms },
		Row: func(e *terrariumpb.Platform) []string {
			shortSha := e.RepoCommit
			if len(e.RepoCommit) > 7 {
				shortSha = e.RepoCommit[:7]
			}
			return []string{e.Id, e.Title, e.RepoUrl, shortSha}
		},
	}

	return f.WriteJsonOrTable(flagOutputFormat == "json")
}
