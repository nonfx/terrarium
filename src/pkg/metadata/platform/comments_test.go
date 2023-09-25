// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDoc(t *testing.T) {
	type args struct {
		filename   string
		endBytePos int
		reverse    bool
	}
	tests := []struct {
		name     string
		args     args
		wantArgs map[string]string
		wantErr  bool
	}{
		{
			name: "read forward N bytes",
			args: args{
				filename:   "./test-component/main.tf",
				endBytePos: 10,
				reverse:    false,
			},
			wantArgs: map[string]string{
				"description": "Auto-gen",
			},
		},
		{
			name: "read forward until EOF",
			args: args{
				filename:   "./test-component/main.tf",
				endBytePos: -1,
				reverse:    false,
			},
			wantArgs: map[string]string{
				"description": "Auto-generated component input definitions.",
				"title":       "Component Inputs",
			},
		},
		{
			name: "read backward",
			args: args{
				filename:   "./test-component/main.tf",
				endBytePos: 340,
				reverse:    true,
			},
			wantArgs: map[string]string{
				"description": "Version of the PostgreSQL engine to use",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotErr := GetDoc(tt.args.filename, tt.args.endBytePos, tt.args.reverse)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.wantArgs, gotArgs)
			}
		})
	}
}

func TestSetValueFromDocIfFound(t *testing.T) {
	type args struct {
		value        *string
		valueTagName string
		fieldDoc     map[string]string
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
	}{
		{
			name: "tag found",
			args: args{
				value:        new(string),
				valueTagName: "moratorium",
				fieldDoc: map[string]string{
					"moratorium": "withdrawal turquoise Incredible",
				},
			},
			wantValue: "withdrawal turquoise Incredible",
		},
		{
			name: "tag not found",
			args: args{
				value:        new(string),
				valueTagName: "optimize",
				fieldDoc: map[string]string{
					"moratorium": "withdrawal turquoise Incredible",
				},
			},
			wantValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetValueFromDocIfFound(tt.args.value, tt.args.valueTagName, tt.args.fieldDoc)
			assert.Equal(t, tt.wantValue, *tt.args.value)
		})
	}
}
