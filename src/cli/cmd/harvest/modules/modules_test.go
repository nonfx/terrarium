// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package modules

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

func createMockVariable() tfValue {
	return &tfconfig.Variable{
		Name:        "dynamic",
		Type:        "Dynamic",
		Description: "alarm",
		Default:     nil,
		Required:    false,
		Sensitive:   true,
		Pos:         tfconfig.SourcePos{},
	}
}

func createMockReference(typeName string, path string) tfconfig.AttributeReference {
	return tfconfig.ResourceAttributeReference{
		Expression:    &hclsyntax.LiteralValueExpr{},
		Module:        &tfconfig.Module{},
		ResourceType:  typeName,
		ResourceName:  "application",
		AttributePath: []string{path},
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

	dbMocks.On("CreateTFModuleAttribute", mock.Anything, mock.Anything).Return(
		func(e *db.TFModuleAttribute) uuid.UUID {
			return uuid.New()
		},
		func(e *db.TFModuleAttribute) error {
			return nil
		})
}

func Test_createAttributeRecord(t *testing.T) {
	type args struct {
		moduleDB         *db.TFModule
		v                tfValue
		varAttributePath string
		res              tfconfig.AttributeReference
	}
	tests := []struct {
		name    string
		args    args
		mockDB  func(*mocks.DB)
		want    *db.TFModuleAttribute
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				moduleDB:         &db.TFModule{},
				v:                createMockVariable(),
				varAttributePath: "Account",
				res:              createMockReference("Albania", "neutral"),
			},
			mockDB: dbMocker,
			want: &db.TFModuleAttribute{
				ModuleAttributeName: "dynamic.Account",
				Description:         "alarm",
				Optional:            true,
				Computed:            false,
			},
			wantErr: false,
		},
		{
			name: "invalid resource kind (error ignored)",
			args: args{
				moduleDB:         &db.TFModule{},
				v:                createMockVariable(),
				varAttributePath: "Account",
				res:              createMockReference("module", "neutral"), // module type
			},
			mockDB:  dbMocker,
			want:    nil,
			wantErr: false,
		},
		{
			name: "invalid resource kind (error ignored)",
			args: args{
				moduleDB:         &db.TFModule{},
				v:                createMockVariable(),
				varAttributePath: "Account",
				res:              createMockReference("var", "neutral"), // var type
			},
			mockDB:  dbMocker,
			want:    nil,
			wantErr: false,
		},
		{
			name: "invalid resource kind (error ignored)",
			args: args{
				moduleDB:         &db.TFModule{},
				v:                createMockVariable(),
				varAttributePath: "Account",
				res:              createMockReference("Fresh", ""), // "" path
			},
			mockDB:  dbMocker,
			want:    nil,
			wantErr: false,
		},
		{
			name: "error from GetTFResourceType (error ignored)",
			args: args{
				moduleDB:         &db.TFModule{},
				v:                createMockVariable(),
				varAttributePath: "Account",
				res:              createMockReference("[fail]", "neutral"),
			},
			mockDB:  dbMocker,
			want:    nil,
			wantErr: false,
		},
		{
			name: "error from GetTFResourceAttribute (error ignored)",
			args: args{
				moduleDB:         &db.TFModule{},
				v:                createMockVariable(),
				varAttributePath: "Account",
				res:              createMockReference("Designer", "[fail]"),
			},
			mockDB:  dbMocker,
			want:    nil,
			wantErr: false,
		},
		{
			name: "error from CreateTFModuleAttribute",
			args: args{
				moduleDB:         &db.TFModule{},
				v:                createMockVariable(),
				varAttributePath: "Account",
				res:              createMockReference("Designer", "payment"),
			},
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("CreateTFModuleAttribute", mock.Anything, mock.Anything).Return(uuid.New(), fmt.Errorf("CreateTFModuleAttribute error")).Once()
				dbMocker(dbMocks)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			if tt.mockDB != nil {
				tt.mockDB(dbMocks)
			}

			resourceTypeByName := make(map[string]*db.TFResourceType)
			got, err := createAttributeRecord(dbMocks, tt.args.moduleDB, tt.args.v, tt.args.varAttributePath, tt.args.res, resourceTypeByName)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAttributeRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
