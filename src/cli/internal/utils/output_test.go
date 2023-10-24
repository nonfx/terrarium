// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/stretchr/testify/assert"
)

func TestOutputFormatter(t *testing.T) {
	tests := []struct {
		name     string
		isJson   bool
		data     *terrariumpb.ListTaxonomyResponse
		expected string
	}{
		{
			name:   "Test JSON Output",
			isJson: true,
			data: &terrariumpb.ListTaxonomyResponse{
				Page: &terrariumpb.Page{Total: 2, Size: 10},
				Taxonomy: []*terrariumpb.Taxonomy{
					{Id: "apple"},
					{Id: "banana"},
				},
			},
			expected: `{"taxonomy":[{"id":"apple", "levels":[]}, {"id":"banana", "levels":[]}], "page":{"size":10, "index":0, "total":2}}`,
		},
		{
			name:   "Test Table Output",
			isJson: false,
			data: &terrariumpb.ListTaxonomyResponse{
				Page: &terrariumpb.Page{Total: 2, Size: 10},
				Taxonomy: []*terrariumpb.Taxonomy{
					{Id: "apple"},
					{Id: "banana"},
				},
			},
			expected: "  #  ITEM    \n  1  apple   \n  2  banana  \n\nPage: 1 of 2 | Page Size: 10\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			f := OutputFormatter[*terrariumpb.ListTaxonomyResponse, *terrariumpb.Taxonomy]{
				Writer:     &buf,
				Data:       tt.data,
				RowHeaders: []string{"Item"},
				Array: func(data *terrariumpb.ListTaxonomyResponse) []*terrariumpb.Taxonomy {
					return data.Taxonomy
				},
				Row: func(item *terrariumpb.Taxonomy) []string {
					return []string{item.Id}
				},
			}

			err := f.WriteJsonOrTable(tt.isJson)
			assert.NoError(t, err)

			actual := buf.String()
			if tt.isJson {
				assert.JSONEq(t, tt.expected, actual)
			} else {
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
