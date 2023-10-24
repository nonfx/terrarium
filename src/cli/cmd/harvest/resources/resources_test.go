// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"fmt"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/tf/schema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_pushProvidersSchemaToDB(t *testing.T) {
	tests := []struct {
		name            string
		providersSchema *schema.ProvidersSchema
		mocks           func(*mocks.DB)

		wantErr           bool
		wantProviderCount int
		wantAllResCount   int
		wantAllAttrCount  int
	}{
		{
			name: "success",
			providersSchema: &schema.ProvidersSchema{
				ProviderSchemas: map[string]schema.ProviderSchema{
					"mock_provider": {
						ResourceSchemas: map[string]schema.SchemaRepresentation{
							"mock_resource": {
								Block: schema.BlockRepresentation{
									Attributes: map[string]schema.AttributeRepresentation{
										"A": {Type: "string"},
									},
									BlockTypes: map[string]schema.BlockTypeRepresentation{
										"X": {
											Block: schema.BlockRepresentation{
												Attributes: map[string]schema.AttributeRepresentation{
													"Y": {Type: "string"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			mocks: func(dbMocks *mocks.DB) {
				dbMocks.On("GetOrCreateTFProvider", mock.Anything).Return(uuid.New(), true, nil).Once()
				dbMocks.On("CreateTFResourceType", mock.Anything).Return(uuid.New(), nil).Once()
				dbMocks.On("CreateTFResourceAttribute", mock.Anything).Return(uuid.New(), nil).Times(3)
			},
			wantProviderCount: 1,
			wantAllResCount:   1,
			wantAllAttrCount:  3,
		},
		{
			name: "panic",
			providersSchema: &schema.ProvidersSchema{
				ProviderSchemas: map[string]schema.ProviderSchema{
					"mock_provider": {},
				},
			},
			mocks: func(dbMocks *mocks.DB) {
				dbMocks.On("GetOrCreateTFProvider", mock.Anything).Return(uuid.Nil, false, fmt.Errorf("mocked error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			tt.mocks(dbMocks)

			providerCount, allResCount, allAttrCount, err := pushProvidersSchemaToDB(tt.providersSchema, dbMocks)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantProviderCount, providerCount)
			assert.Equal(t, tt.wantAllResCount, allResCount)
			assert.Equal(t, tt.wantAllAttrCount, allAttrCount)

			dbMocks.AssertExpectations(t)

		})
	}
}

func Test_loadProvidersSchema(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.ProvidersSchema
		wantErr bool
	}{
		{
			name: "success",
			args: args{"./testdata/example_schema.json"},
			want: &schema.ProvidersSchema{
				ProviderSchemas: map[string]schema.ProviderSchema{
					"mock_provider": {
						ResourceSchemas: map[string]schema.SchemaRepresentation{
							"mock_resource": {
								Block: schema.BlockRepresentation{
									Attributes: map[string]schema.AttributeRepresentation{
										"A": {
											Description: "a",
											Type:        "string",
											Computed:    true,
										},
									},
									BlockTypes: map[string]schema.BlockTypeRepresentation{
										"X": {
											Block: schema.BlockRepresentation{
												Attributes: map[string]schema.AttributeRepresentation{
													"Y": {
														Description: "y",
														Type:        "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "fail read",
			args:    args{"./invalid_file_path"},
			wantErr: true,
		},
		{
			name:    "fail unmarshal",
			args:    args{"./testdata/invalid_schema.txt"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadProvidersSchema(tt.args.filename)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
