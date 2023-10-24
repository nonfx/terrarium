// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const (
	simpleExample = "  A SIMPLE    EXAMPLE   OF    LIST    OUTPUT\n"
	skipActual    = "  A SIMPLE    NOTEXAMPLE   OF    LIST    OUTPUT\n"
)

func Test_splitListOutput(t *testing.T) {
	type args struct {
		output string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test 1",
			args: args{output: simpleExample},
			want: []string{"A SIMPLE", "EXAMPLE", "OF", "LIST", "OUTPUT"},
		},
		{
			name: "Empty string",
			args: args{""},
			want: []string{},
		},
		{
			name: "All whitespace",
			args: args{"   \n   \n    \t"},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitListOutput(tt.args.output); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitListOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateListOutput(t *testing.T) {
	type args struct {
		t          *testing.T
		expected   string
		actual     string
		skipFields []int
	}
	tests := []struct {
		name    string
		args    args
		isValid bool
	}{
		{
			name: "Success no skip",
			args: args{
				t:          new(testing.T),
				expected:   simpleExample,
				actual:     simpleExample,
				skipFields: []int{},
			},
			isValid: true,
		},
		{
			name: "Success with skip",
			args: args{
				t:          new(testing.T),
				expected:   simpleExample,
				actual:     skipActual,
				skipFields: []int{2},
			},
			isValid: true,
		},
		{
			name: "Failure no skip",
			args: args{
				t:          new(testing.T),
				expected:   simpleExample,
				actual:     "Not a simple example",
				skipFields: []int{},
			},
			isValid: false,
		},
		{
			name: "Failure with skip",
			args: args{
				t:          new(testing.T),
				expected:   simpleExample,
				actual:     skipActual,
				skipFields: []int{1},
			},
			isValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, ValidateListOutput(tt.args.t, tt.args.expected, tt.args.actual, tt.args.skipFields))
		})
	}
}
