// Copyright (c) CloudCover
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

func TestService_ListComponents(t *testing.T) {

	TestCases[terrariumpb.ListComponentsRequest, terrariumpb.ListComponentsResponse]{
		{
			name: "success",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListComponentsRequest, terrariumpb.ListComponentsResponse]) {
				tc.mockDB.On("QueryPlatformComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(db.PlatformComponents{
					{
						Dependency: db.Dependency{Title: "mocked"},
					},
				}, nil)
			},
			req: &terrariumpb.ListComponentsRequest{
				Page: &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
			},
			wantResp: &terrariumpb.ListComponentsResponse{
				Page: &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
				Components: []*terrariumpb.Component{
					{
						Id:            uuid.Nil.String(),
						InterfaceUuid: uuid.Nil.String(),
						Title:         "mocked",
					},
				},
			},
		},
		{
			name: "db query error",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListComponentsRequest, terrariumpb.ListComponentsResponse]) {
				tc.mockDB.On("QueryPlatformComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("mocked err"))
			},
			req: &terrariumpb.ListComponentsRequest{
				Page: &terrariumpb.Page{Size: 10, Index: 2},
			},
			wantErr: "error running database query: mocked err",
		},
	}.Run(t, func(s *Service) func(context.Context, *terrariumpb.ListComponentsRequest) (*terrariumpb.ListComponentsResponse, error) {
		return s.ListComponents
	})
}
