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

func TestService_ListTaxonomy(t *testing.T) {

	TestCases[terrariumpb.ListTaxonomyRequest, terrariumpb.ListTaxonomyResponse]{
		{
			name: "success",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListTaxonomyRequest, terrariumpb.ListTaxonomyResponse]) {
				tc.mockDB.On("QueryTaxonomies", mock.Anything, mock.Anything).Return(db.Taxonomies{
					*db.TaxonomyFromLevels("mocked-l1", "l2", "l3"),
				}, nil)
			},
			req: &terrariumpb.ListTaxonomyRequest{
				Taxonomy: []string{"mocked-l1"},
				Page:     &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
			},
			wantResp: &terrariumpb.ListTaxonomyResponse{
				Page: &terrariumpb.Page{Size: 10, Index: 2, Total: 1},
				Taxonomy: []*terrariumpb.Taxonomy{
					{
						Id:     uuid.Nil.String(),
						Levels: []string{"mocked-l1", "l2", "l3"},
					},
				},
			},
		},
		{
			name: "db query error",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListTaxonomyRequest, terrariumpb.ListTaxonomyResponse]) {
				tc.mockDB.On("QueryTaxonomies", mock.Anything, mock.Anything).Return(nil, errors.New("mocked err"))
			},
			req: &terrariumpb.ListTaxonomyRequest{
				Page: &terrariumpb.Page{Size: 10, Index: 2},
			},
			wantErr: "error running database query: mocked err",
		},
	}.Run(t, func(s *Service) func(context.Context, *terrariumpb.ListTaxonomyRequest) (*terrariumpb.ListTaxonomyResponse, error) {
		return s.ListTaxonomy
	})
}
