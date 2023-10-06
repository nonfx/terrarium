// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"io"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/transporthelper"
	"github.com/olekukonko/tablewriter"
	"github.com/rotisserie/eris"
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

// OutputFormatter is a generic struct that takes two types T1 and T2.
// It is used for formatting output data either as JSON or as a table.
type OutputFormatter[T1, T2 any] struct {
	Writer     io.Writer         // Writer to write the output to
	Data       T1                // Data to be formatted
	RowHeaders []string          // Headers for table rows
	Array      func(T1) []T2     // Function to convert Data into an array of T2
	Row        func(T2) []string // Function to convert each T2 into a row of strings
}

// WriteJson formats the Data as JSON and writes it to the Writer.
func (f OutputFormatter[T1, T2]) WriteJson() error {
	b, err := transporthelper.CreateJSONBodyMarshaler().Marshaler.Marshal(f.Data)
	if err != nil {
		return eris.Wrap(err, "error formatting output to json")
	}

	fmt.Fprintf(f.Writer, "%s\n", b)
	return nil
}

// WriteTable formats the Data as a table and writes it to the Writer.
func (f OutputFormatter[T1, T2]) WriteTable() error {
	table := OutFormatForList(f.Writer)
	table.SetHeader(append([]string{"#"}, f.RowHeaders...))

	// Loop through each item in the array and append it to the table
	for i, res := range f.Array(f.Data) {
		row := f.Row(res)
		row = append([]string{fmt.Sprintf("%d", i+1)}, row...)
		table.Append(row)
	}

	table.Render()

	// If Data has a GetPage method, display pagination information
	if pageGetter, ok := ((interface{})(f.Data)).(interface{ GetPage() *terrariumpb.Page }); ok {
		p := pageGetter.GetPage()
		fmt.Fprintf(f.Writer, "\nPage: %d of %d | Page Size: %d\n", p.Index+1, p.Total, p.Size)
	}

	return nil
}

// WriteJsonOrTable decides whether to write Data as JSON or as a table based on the isJson flag.
func (f OutputFormatter[T1, T2]) WriteJsonOrTable(isJson bool) error {
	if isJson {
		return f.WriteJson()
	}

	return f.WriteTable()
}
