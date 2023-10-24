// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestIsBool(t *testing.T) {
	type args struct {
		expr hcl.Expression
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "bool expression",
			args: args{
				expr: &hclsyntax.BinaryOpExpr{},
			},
			want: true,
		},
		{
			name: "bool constant",
			args: args{
				expr: &hclsyntax.LiteralValueExpr{
					Val: cty.BoolVal(false),
				},
			},
			want: true,
		},
		{
			name: "non-bool constant",
			args: args{
				expr: &hclsyntax.LiteralValueExpr{
					Val: cty.StringVal("n{25)uOaI_"),
				},
			},
			want: false,
		},
		{
			name: "bool function",
			args: args{
				expr: &hclsyntax.FunctionCallExpr{
					Name: "anytrue",
				},
			},
			want: true,
		},
		{
			name: "non-bool function",
			args: args{
				expr: &hclsyntax.FunctionCallExpr{
					Name: "$g8C$eQ1qT",
				},
			},
			want: false,
		},
		{
			name: "non-bool expression",
			args: args{
				expr: &hclsyntax.ObjectConsExpr{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBool(tt.args.expr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsObject(t *testing.T) {
	type args struct {
		expr hcl.Expression
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "object literal expression",
			args: args{
				expr: &hclsyntax.ObjectConsExpr{},
			},
			want: true,
		},
		{
			name: "other expression",
			args: args{
				expr: &hclsyntax.BinaryOpExpr{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsObject(tt.args.expr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsCollection(t *testing.T) {
	type args struct {
		expr hcl.Expression
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "dynamic collection expression",
			args: args{
				expr: &hclsyntax.ForExpr{},
			},
			want: true,
		},
		{
			name: "object literal expression",
			args: args{
				expr: &hclsyntax.ObjectConsExpr{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCollection(tt.args.expr)
			assert.Equal(t, tt.want, got)
		})
	}
}
