// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package lint

import (
	"testing"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func Test_lintPlatform(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid platform - default not in enum",
			args: args{
				dir: "testdata/invalid-terraform-1",
			},
			wantErr: true,
		},
		{
			name: "valid platform",
			args: args{
				dir: "testdata/valid-terraform-1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lintPlatform(tt.args.dir)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_validatePlatformTerraform(t *testing.T) {
	mockStringExpr := hcl.StaticExpr(cty.StringVal("Human"), hcl.Range{
		Filename: "quality_kroon_generate.twds",
		Start: hcl.Pos{
			Line:   1,
			Column: 2,
			Byte:   3,
		},
		End: hcl.Pos{
			Line:   1,
			Column: 2,
			Byte:   3,
		},
	})
	mockBoolExpr := hclsyntax.LiteralValueExpr{
		Val:      cty.BoolVal(true),
		SrcRange: hcl.Range{},
	}
	mockObjectExpr := hclsyntax.ObjectConsExpr{
		Items:     []hclsyntax.ObjectConsItem{},
		SrcRange:  hcl.Range{},
		OpenRange: hcl.Range{},
	}
	type args struct {
		module *tfconfig.Module
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "inputs not object",
			args: args{
				module: &tfconfig.Module{
					Path:     "",
					Metadata: &tfconfig.Metadata{},
					Locals: map[string]*tfconfig.Local{
						"tr_component_postgres": {
							Name:       "tr_component_postgres",
							Expression: mockStringExpr,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "taxon switch not bool",
			args: args{
				module: &tfconfig.Module{
					Path:     "",
					Metadata: &tfconfig.Metadata{},
					Locals: map[string]*tfconfig.Local{
						"tr_taxon_db_enabled": {
							Name:       "tr_taxon_db_enabled",
							Expression: mockStringExpr,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "module switch not bool",
			args: args{
				module: &tfconfig.Module{
					Path:     "",
					Metadata: &tfconfig.Metadata{},
					Locals: map[string]*tfconfig.Local{
						"tr_component_postgres": {
							Name:       "tr_component_postgres",
							Expression: &mockObjectExpr,
						},
						"tr_component_postgres_enabled": {
							Name:       "tr_component_postgres_enabled",
							Expression: mockStringExpr,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing module call",
			args: args{
				module: &tfconfig.Module{
					Path:     "",
					Metadata: &tfconfig.Metadata{},
					Locals: map[string]*tfconfig.Local{
						"tr_component_postgres": {
							Name:       "tr_component_postgres",
							Expression: &mockObjectExpr,
						},
						"tr_component_postgres_enabled": {
							Name:       "tr_component_postgres_enabled",
							Expression: &mockBoolExpr,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "output not collection",
			args: args{
				module: &tfconfig.Module{
					Path:     "",
					Metadata: &tfconfig.Metadata{},
					Locals: map[string]*tfconfig.Local{
						"tr_component_postgres": {
							Name:       "tr_component_postgres",
							Expression: &mockObjectExpr,
						},
						"tr_component_postgres_enabled": {
							Name:       "tr_component_postgres_enabled",
							Expression: &mockBoolExpr,
						},
					},
					ModuleCalls: map[string]*tfconfig.ModuleCall{
						"tr_component_postgres": {},
					},
					Outputs: map[string]*tfconfig.Output{
						"tr_component_postgres_host": {
							Value: tfconfig.ResourceAttributeReference{
								Expression: mockStringExpr,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "output not collection",
			args: args{
				module: &tfconfig.Module{
					Path:     "",
					Metadata: &tfconfig.Metadata{},
					Locals: map[string]*tfconfig.Local{
						"tr_component_postgres": {
							Name:       "tr_component_postgres",
							Expression: &mockObjectExpr,
						},
						"tr_component_postgres_enabled": {
							Name:       "tr_component_postgres_enabled",
							Expression: &mockBoolExpr,
						},
					},
					ModuleCalls: map[string]*tfconfig.ModuleCall{
						"tr_component_postgres": {},
					},
					Outputs: map[string]*tfconfig.Output{
						"tr_component_postgres_host": {
							Value: tfconfig.ResourceAttributeReference{
								Expression: &hclsyntax.ForExpr{},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePlatformTerraform(tt.args.module); (err != nil) != tt.wantErr {
				t.Errorf("validatePlatformTerraform() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
