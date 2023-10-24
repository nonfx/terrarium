// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"
	"errors"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestService_ListDependencies(t *testing.T) {

	TestCases[terrariumpb.ListDependenciesRequest, terrariumpb.ListDependenciesResponse]{
		{
			name: "success",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListDependenciesRequest, terrariumpb.ListDependenciesResponse]) {
				tc.mockDB.On("QueryDependencies", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(db.Dependencies{
					{Title: "mocked"},
				}, nil)
			},
			req: &terrariumpb.ListDependenciesRequest{
				Page: &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
			},
			wantResp: &terrariumpb.ListDependenciesResponse{
				Page: &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
				Dependencies: []*terrariumpb.Dependency{
					{
						Id:    uuid.Nil.String(),
						Title: "mocked",
					},
				},
			},
		},
		{
			name: "db query error",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListDependenciesRequest, terrariumpb.ListDependenciesResponse]) {
				tc.mockDB.On("QueryDependencies", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("mocked err"))
			},
			req: &terrariumpb.ListDependenciesRequest{
				Page: &terrariumpb.Page{Size: 10, Index: 2},
			},
			wantErr: "error running database query: mocked err",
		},
	}.Run(t, func(s *Service) func(context.Context, *terrariumpb.ListDependenciesRequest) (*terrariumpb.ListDependenciesResponse, error) {
		return s.ListDependencies
	})
}
