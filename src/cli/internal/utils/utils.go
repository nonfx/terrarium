package utils

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// OutFormatForList This configures the table format for the List commands
func OutFormatForList(out io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(out)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	return table
}
