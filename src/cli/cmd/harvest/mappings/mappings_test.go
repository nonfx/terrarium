// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package mappings

import (
	"fmt"
	"testing"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type resCfg struct {
	ParentType          string
	ParentAttributePath string
	RefType             string
	RefAttributePath    string
}

func (c resCfg) GetParentType() string {
	return valueOrRandom(c.ParentType)
}

func (c resCfg) GetParentAttributePath() string {
	return valueOrRandom(c.ParentAttributePath)
}

func (c resCfg) GetRefType() string {
	return valueOrRandom(c.RefType)
}

func (c resCfg) GetRefAttributePath() string {
	return valueOrRandom(c.RefAttributePath)
}

func valueOrRandom(value string) string {
	if value != "" {
		return value
	}
	return uuid.New().String()
}

type cfg struct {
	Managed resCfg
	Data    resCfg
}

func createTestModule(config cfg, subMod *tfconfig.Module) *tfconfig.Module {
	return &tfconfig.Module{
		Path: "Rustic",
		ManagedResources: map[string]*tfconfig.Resource{
			"Mouse": {
				Type: config.Managed.GetParentType(),
				Name: "Mouse",
				Provider: tfconfig.ProviderRef{
					Name: "Loan",
				},
				References: map[string][]tfconfig.AttributeReference{
					config.Managed.GetParentAttributePath(): {
						tfconfig.ResourceAttributeReference{
							Expression:    &hclsyntax.LiteralValueExpr{},
							Module:        &tfconfig.Module{},
							ResourceType:  config.Managed.GetRefType(),
							ResourceName:  "Awesome",
							AttributePath: []string{config.Managed.GetRefAttributePath()},
						},
					},
				},
			},
			"Metrics": {
				Type: "Dynamic",
				Name: "Metrics",
				Provider: tfconfig.ProviderRef{
					Name: "Rubber",
				},
				References: map[string][]tfconfig.AttributeReference{
					"Hawaii": {
						tfconfig.ResourceAttributeReference{
							Expression:    &hclsyntax.LiteralValueExpr{},
							Module:        &tfconfig.Module{},
							ResourceType:  "next",
							ResourceName:  "Loan",
							AttributePath: []string{"e-services"},
						},
					},
				},
			},
		},
		DataResources: map[string]*tfconfig.Resource{
			"Buckinghamshire": {
				Type: config.Data.GetParentType(),
				Name: "Buckinghamshire",
				Provider: tfconfig.ProviderRef{
					Name: "Buckinghamshire",
				},
				References: map[string][]tfconfig.AttributeReference{
					config.Data.GetParentAttributePath(): {
						tfconfig.ResourceAttributeReference{
							Expression:    &hclsyntax.LiteralValueExpr{},
							Module:        &tfconfig.Module{},
							ResourceType:  config.Data.GetRefType(),
							ResourceName:  "Cambridgeshire",
							AttributePath: []string{config.Data.GetRefAttributePath()},
						},
					},
				},
			},
			"Loan": {
				Type: "Sausages",
				Name: "Loan",
				Provider: tfconfig.ProviderRef{
					Name: "intermediate",
				},
				References: map[string][]tfconfig.AttributeReference{
					"index": {
						tfconfig.ResourceAttributeReference{
							Expression:    &hclsyntax.LiteralValueExpr{},
							Module:        &tfconfig.Module{},
							ResourceType:  "var", // ignored type
							ResourceName:  "Cambridgeshire",
							AttributePath: []string{"Maine"},
						},
					},
				},
			},
		},
		ModuleCalls: map[string]*tfconfig.ModuleCall{
			"SAS": {
				Name:    "ROI",
				Source:  "solutions",
				Version: "input",
				Module:  subMod,
			},
		},
	}
}

func dbMocker(dbMocks *mocks.DB) {
	dbMocks.On("GetTFResourceType", mock.Anything, mock.Anything).Return(func(e *db.TFResourceType, where *db.TFResourceType) error {
		if where.ResourceType == "[fail]" {
			return fmt.Errorf("error from GetTFResourceType")
		}
		*e = *where
		return nil
	})
	dbMocks.On("GetTFResourceAttribute", mock.Anything, mock.Anything).Return(func(e *db.TFResourceAttribute, where *db.TFResourceAttribute) error {
		if where.AttributePath == "[fail]" {
			return fmt.Errorf("error from GetTFResourceAttribute")
		}
		*e = *where
		return nil
	})

	dbMocks.On("CreateTFResourceAttributesMapping", mock.Anything, mock.Anything).Return(
		func(e *db.TFResourceAttributesMapping) uuid.UUID {
			return uuid.New()
		},
		func(e *db.TFResourceAttributesMapping) error {
			return nil
		})
}

func Test_createMappingsForModule(t *testing.T) {
	type args struct {
		config *tfconfig.Module
	}
	tests := []struct {
		name              string
		args              args
		mockDB            func(*mocks.DB)
		wantResourceCount int
		wantMappingsCount int
		wantErr           bool
	}{
		{
			name: "success",
			args: args{
				config: createTestModule(cfg{}, createTestModule(cfg{}, nil)),
			},
			mockDB:            dbMocker,
			wantResourceCount: 8,
			wantMappingsCount: 6,
			wantErr:           false,
		},
		{
			name: "error from dest GetTFResourceType (error ignored)",
			args: args{
				config: createTestModule(cfg{
					Managed: resCfg{
						ParentType: "[fail]",
					},
				}, createTestModule(cfg{}, nil)),
			},
			mockDB:            dbMocker,
			wantResourceCount: 8,
			wantMappingsCount: 5,
			wantErr:           false,
		},
		{
			name: "error from src GetTFResourceType (error ignored)",
			args: args{
				config: createTestModule(cfg{
					Managed: resCfg{
						RefType: "[fail]",
					},
				}, createTestModule(cfg{}, nil)),
			},
			mockDB:            dbMocker,
			wantResourceCount: 8,
			wantMappingsCount: 5,
			wantErr:           false,
		},
		{
			name: "error from dest GetTFResourceAttribute",
			args: args{
				config: createTestModule(cfg{
					Managed: resCfg{
						ParentAttributePath: "[fail]",
					},
					Data: resCfg{},
				}, createTestModule(cfg{}, nil)),
			},
			mockDB:  dbMocker,
			wantErr: true,
		},
		{
			name: "error from src GetTFResourceAttribute (error ignored)",
			args: args{
				config: createTestModule(cfg{
					Managed: resCfg{
						RefAttributePath: "[fail]",
					},
					Data: resCfg{},
				}, createTestModule(cfg{}, nil)),
			},
			mockDB:            dbMocker,
			wantResourceCount: 8,
			wantMappingsCount: 5,
			wantErr:           false,
		},
		{
			name: "error from CreateTFResourceAttributesMapping",
			args: args{
				config: createTestModule(cfg{}, createTestModule(cfg{}, nil)),
			},
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("CreateTFResourceAttributesMapping", mock.Anything, mock.Anything).Return(uuid.New(), fmt.Errorf("CreateTFResourceAttributesMapping error")).Once()
				dbMocker(dbMocks)
			},
			wantErr: true,
		},
		{
			name: "error from DataResource",
			args: args{
				config: createTestModule(cfg{
					Data: resCfg{
						ParentAttributePath: "[fail]",
					},
				}, createTestModule(cfg{}, nil)),
			},
			mockDB:  dbMocker,
			wantErr: true,
		},
		{
			name: "error from sub-module",
			args: args{
				config: createTestModule(cfg{}, createTestModule(cfg{
					Data: resCfg{
						ParentAttributePath: "[fail]",
					},
				}, nil)),
			},
			mockDB:  dbMocker,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			if tt.mockDB != nil {
				tt.mockDB(dbMocks)
			}

			gotMappings, gotResourceCount, err := createMappingsForModule(dbMocks, tt.args.config, make(map[string]*db.TFResourceType))
			if (err != nil) != tt.wantErr {
				t.Errorf("createMappingsForModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantResourceCount, gotResourceCount)
			assert.Len(t, gotMappings, tt.wantMappingsCount)
		})
	}
}
