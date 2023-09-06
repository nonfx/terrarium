// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package modules

import (
	"fmt"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/cli/internal/utils"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/transporthelper"
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	flagPopulateMappings bool
	flagSearchText       string
	flagOutputFormat     string
	flagPageSize         int32
	flagPageIndex        int32
	flagNamespaces       []string
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "modules",
		Short: "List modules matching the source pattern",
		Long:  "command to list mathcing modules as per the filter chosen. provides variety of filters to list desired modules",
		RunE:  listModules,
	}

	cmd.Flags().StringVarP(&flagSearchText, "searchText", "s", "", "optional search text")
	cmd.Flags().BoolVarP(&flagPopulateMappings, "populateMappings", "p", false, "A boolean flag to populate mappings")
	cmd.Flags().Int32Var(&flagPageSize, "pageSize", 100, "page size flag")
	cmd.Flags().Int32Var(&flagPageIndex, "pageIndex", 0, "page index flag")
	cmd.Flags().StringVarP(&flagOutputFormat, "output", "o", "table", "Output format (json or table)")
	cmd.Flags().StringSliceVarP(&flagNamespaces, "namespaces", "n", []string{}, "namespaces filter - farm repo will always be included")

	return cmd
}

func listModules(cmd *cobra.Command, args []string) error {
	g, err := config.DBConnect()
	if err != nil {
		return err
	}

	page := &terrariumpb.Page{
		Size:  flagPageSize,
		Index: flagPageIndex,
	}
	flagNamespaces = append(flagNamespaces, config.FarmDefault())
	result, err := g.QueryTFModules(
		db.ModuleSearchFilter(flagSearchText),
		db.PopulateModuleMappingsFilter(flagPopulateMappings),
		db.PaginateGlobalFilter(page.Size, page.Index, &page.Total),
		db.ModuleNamespaceFilter(flagNamespaces),
	)
	if err != nil {
		return err
	}

	pbRes := &terrariumpb.ListModulesResponse{
		Page:    page,
		Modules: result.ToProto(),
	}

	if flagOutputFormat == "json" {
		b, err := transporthelper.CreateJSONBodyMarshaler().Marshaler.Marshal(pbRes)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), string(b))

	} else {
		table := utils.OutFormatForList(cmd.OutOrStdout())
		table.SetHeader([]string{"ID", "Module Name", "Source", "Version", "Namespace"})
		for _, res := range result {
			outputLine := []string{res.ID.String(), res.ModuleName, res.Source, string(res.Version), res.Namespace}
			table.Append(outputLine)
		}
		table.Render()
	}
	return nil
}
