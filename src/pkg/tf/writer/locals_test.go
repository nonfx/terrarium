// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package writer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToTFLocals(t *testing.T) {
	tests := []struct {
		input   map[string]interface{}
		output  string
		wantErr bool
	}{
		{
			input: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
			output: "locals {\n  key1 = \"value1\"\n  key2 = 42\n}\n",
		},
		{
			input: map[string]interface{}{
				"enabled": true,
				"rate":    0.5,
			},
			output: "locals {\n  enabled = true\n  rate    = 0.5\n}\n",
		},
		{
			input: map[string]interface{}{
				"key1": map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": "val1",
					},
				},
			},
			output: "locals {\n  key1 = {\n    level1 = {\n      level2 = \"val1\"\n    }\n  }\n}\n",
		},
		{
			input: map[string]interface{}{
				"key1": []interface{}{
					1, 2, 3,
				},
			},
			output: "locals {\n  key1 = [1, 2, 3]\n}\n",
		},
		{
			input: map[string]interface{}{
				"key1": struct{}{},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		err := WriteLocals(test.input, &buf)
		assert.Equal(t, err != nil, test.wantErr, "want error=%v, got: %v", test.wantErr, err)
		if !test.wantErr {
			assert.Equal(t, test.output, buf.String())
		}
	}
}
