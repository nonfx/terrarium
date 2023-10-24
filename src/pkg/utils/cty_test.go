// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestToCtyValue(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    cty.Value
		wantErr bool
	}{
		{
			name:  "string",
			input: "test",
			want:  cty.StringVal("test"),
		},
		{
			name:  "int",
			input: 42,
			want:  cty.NumberIntVal(42),
		},
		{
			name:  "float64",
			input: 42.5,
			want:  cty.NumberFloatVal(42.5),
		},
		{
			name:  "bool",
			input: true,
			want:  cty.BoolVal(true),
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"key": map[string]interface{}{
					"nestedKey": "value",
				},
			},
			want: cty.ObjectVal(map[string]cty.Value{
				"key": cty.ObjectVal(map[string]cty.Value{
					"nestedKey": cty.StringVal("value"),
				}),
			}),
		},
		{
			name:  "array",
			input: []interface{}{1, "2", 3},
			want:  cty.TupleVal([]cty.Value{cty.NumberIntVal(1), cty.StringVal("2"), cty.NumberIntVal(3)}),
		},
		{
			name:    "unsupported type",
			input:   struct{}{},
			want:    cty.NilVal,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToCtyValue(tt.input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
