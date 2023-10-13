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

func TestService_ListPlatforms(t *testing.T) {

	TestCases[terrariumpb.ListPlatformsRequest, terrariumpb.ListPlatformsResponse]{
		{
			name: "success",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListPlatformsRequest, terrariumpb.ListPlatformsResponse]) {
				tc.mockDB.On("QueryPlatforms", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(db.Platforms{
					{Title: "mocked-pf1"},
				}, nil)
			},
			req: &terrariumpb.ListPlatformsRequest{
				Page: &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
			},
			wantResp: &terrariumpb.ListPlatformsResponse{
				Page: &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
				Platforms: []*terrariumpb.Platform{
					{
						Id:    uuid.Nil.String(),
						Title: "mocked-pf1",
					},
				},
			},
		},
		{
			name: "db query error",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListPlatformsRequest, terrariumpb.ListPlatformsResponse]) {
				tc.mockDB.On("QueryPlatforms", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("mocked err"))
			},
			req: &terrariumpb.ListPlatformsRequest{
				Page: &terrariumpb.Page{Size: 10, Index: 2},
			},
			wantErr: "error running database query: mocked err",
		},
	}.Run(t, func(s *Service) func(context.Context, *terrariumpb.ListPlatformsRequest) (*terrariumpb.ListPlatformsResponse, error) {
		return s.ListPlatforms
	})
}
