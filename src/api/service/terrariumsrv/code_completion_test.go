// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"
	_ "embed"
	"errors"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
)

//go:embed hcl_go_tmpl/module.test1.tf.dat
var moduleTest1 string

//go:embed hcl_go_tmpl/module.test2.tf.dat
var moduleTest2 string

func TestService_CodeCompletion(t *testing.T) {
	mockUuid1 := uuid.New()
	mockUuid2 := uuid.New()
	mockUuid3 := uuid.New()

	TestCases[terrariumpb.CompletionRequest, terrariumpb.CompletionResponse]{
		{
			name: "Successful basic usage",
			req: &terrariumpb.CompletionRequest{
				Modules: []string{mockUuid1.String()},
			},
			wantResp: &terrariumpb.CompletionResponse{
				Suggestions: []string{moduleTest1},
			},
			preCall: func(t *testing.T, tc TestCase[terrariumpb.CompletionRequest, terrariumpb.CompletionResponse]) {
				tc.mockDB.On("FindOutputMappingsByModuleID", mockUuid1).Return(db.TFModules{
					{
						Model:      db.Model{ID: mockUuid1},
						ModuleName: "mock_module_A",
						Source:     "tf/mock_module_A",
						Attributes: []db.TFModuleAttribute{
							{
								ModuleAttributeName: "module_attr_X",
								Optional:            true,
								Computed:            false,
								ResourceAttribute: &db.TFResourceAttribute{
									AttributePath: "res_attr_X",
									OutputMappings: []db.TFResourceAttributesMapping{
										{
											OutputAttribute: db.TFResourceAttribute{
												AttributePath: "res_attr_Y",
												RelatedModuleAttrs: []db.TFModuleAttribute{
													{
														ModuleID:            mockUuid2,
														ModuleAttributeName: "module_attr_Y",
														Module: &db.TFModule{
															Model:      db.Model{ID: mockUuid2},
															ModuleName: "mock_module_B",
															Source:     "tf/mock_module_B",
															Attributes: []db.TFModuleAttribute{
																{
																	ModuleAttributeName: "module_attr_Y",
																	Optional:            false, // required
																	Computed:            true,
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
				}, nil)
				tc.mockDB.On("FindOutputMappingsByModuleID", mockUuid2).Return(db.TFModules{
					{
						Model:      db.Model{ID: mockUuid2},
						ModuleName: "mock_module_B",
						Source:     "tf/mock_module_B",
						Attributes: []db.TFModuleAttribute{
							{
								ModuleAttributeName: "module_attr_Y",
								Computed:            true,
							},
						},
					},
				}, nil)
			},
			postCheck: func(t *testing.T, tc TestCase[terrariumpb.CompletionRequest, terrariumpb.CompletionResponse], resp *terrariumpb.CompletionResponse) {
				// // update test data with returned string
				// writeToFile(t, "hcl_go_tmpl/module.test1.tf.dat", resp.Suggestions[0])
			},
		},
		{
			name: "Successful with context",
			req: &terrariumpb.CompletionRequest{
				CodeContext: moduleTest1,
				Modules:     []string{mockUuid3.String()},
			},
			wantResp: &terrariumpb.CompletionResponse{
				Suggestions: []string{moduleTest2},
			},
			preCall: func(t *testing.T, tc TestCase[terrariumpb.CompletionRequest, terrariumpb.CompletionResponse]) {
				tc.mockDB.On("FindOutputMappingsByModuleID", mockUuid3).Return(db.TFModules{
					{
						Model:      db.Model{ID: mockUuid3},
						ModuleName: "mock_module_C",
						Source:     "tf/mock_module_C",
						Version:    "1.0.0",
						Attributes: []db.TFModuleAttribute{
							{
								ModuleAttributeName: "module_attr_X",
								Optional:            true,
								Computed:            false,
								ResourceAttribute: &db.TFResourceAttribute{
									AttributePath: "res_attr_X",
									OutputMappings: []db.TFResourceAttributesMapping{
										{
											OutputAttribute: db.TFResourceAttribute{
												AttributePath: "res_attr_Y",
												RelatedModuleAttrs: []db.TFModuleAttribute{
													{
														ModuleID:            mockUuid2,
														ModuleAttributeName: "module_attr_Y",
														Module: &db.TFModule{
															Model:      db.Model{ID: mockUuid2},
															ModuleName: "mock_module_B",
															Source:     "tf/mock_module_B",
															Attributes: []db.TFModuleAttribute{
																{
																	ModuleAttributeName: "module_attr_Y",
																	Optional:            false, // required
																	Computed:            true,
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
				}, nil)
				tc.mockDB.On("FindOutputMappingsByModuleID", mockUuid2).Return(db.TFModules{
					{
						Model:      db.Model{ID: mockUuid2},
						ModuleName: "mock_module_B",
						Source:     "tf/mock_module_B",
						Attributes: []db.TFModuleAttribute{
							{
								ModuleAttributeName: "module_attr_Y",
								Computed:            true,
							},
						},
					},
				}, nil)
			},
			postCheck: func(t *testing.T, tc TestCase[terrariumpb.CompletionRequest, terrariumpb.CompletionResponse], resp *terrariumpb.CompletionResponse) {
				// // update test data with returned string
				// writeToFile(t, "hcl_go_tmpl/module.test2.tf.dat", resp.Suggestions[0])
			},
		},
		{
			name: "fail to fetch modules",
			req: &terrariumpb.CompletionRequest{
				CodeContext: moduleTest1,
				Modules:     []string{mockUuid3.String()},
			},
			wantErr: errors.New("MOCKED-ERR failed to retrieve module dependencies"),
			preCall: func(t *testing.T, tc TestCase[terrariumpb.CompletionRequest, terrariumpb.CompletionResponse]) {
				tc.mockDB.On("FindOutputMappingsByModuleID", mockUuid3).Return(nil, tc.wantErr)
			},
		},
	}.Run(t, func(s *Service) func(context.Context, *terrariumpb.CompletionRequest) (*terrariumpb.CompletionResponse, error) {
		return s.CodeCompletion
	})
}
