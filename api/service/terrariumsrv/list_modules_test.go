package terrariumsrv

import (
	"context"
	"errors"
	"testing"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	"github.com/google/uuid"
)

func TestService_ListModules(t *testing.T) {
	mockUuid1 := uuid.New()

	TestCases[terrariumpb.ListModulesRequest, terrariumpb.ListModulesResponse]{
		{
			name: "Successful retrieval of modules",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListModulesRequest, terrariumpb.ListModulesResponse]) {
				tc.mockDB.On("ListTFModule", tc.req.Search, int(tc.req.Page.Size), int(tc.req.Page.Index*tc.req.Page.Size)).
					Return(db.TFModules{
						{
							Model: db.Model{
								ID: mockUuid1,
							},
							ModuleName:  "Rds",
							Version:     "1",
							Source:      "/rds",
							Description: "",
						},
					}, int64(1), nil)
			},
			req: &terrariumpb.ListModulesRequest{
				Search: "search query",
				Page: &terrariumpb.Page{
					Size:  10,
					Index: 2,
				},
			},
			wantResp: &terrariumpb.ListModulesResponse{
				Page: &terrariumpb.Page{
					Size:  10,
					Index: 2,
					Total: 1,
				},
				Modules: []*terrariumpb.Module{
					{
						Id:         mockUuid1.String(),
						TaxonomyId: uuid.Nil.String(),
						ModuleName: "Rds",
						Version:    "1",
						Source:     "/rds",
					},
				},
			},
		},
		{
			name: "Error retrieving modules",
			preCall: func(t *testing.T, tc TestCase[terrariumpb.ListModulesRequest, terrariumpb.ListModulesResponse]) {
				tc.mockDB.On("ListTFModule", tc.req.Search, int(tc.req.Page.Size), int(tc.req.Page.Index*tc.req.Page.Size)).
					Return(nil, int64(0), tc.wantErr)
			},
			req: &terrariumpb.ListModulesRequest{
				Search: "search query",
				Page: &terrariumpb.Page{
					Size:  10,
					Index: 2,
				},
			},
			wantErr: errors.New("MOCKED-ERR failed to retrieve modules"),
		},
	}.Run(t, func(s *Service) func(context.Context, *terrariumpb.ListModulesRequest) (*terrariumpb.ListModulesResponse, error) {
		return s.ListModules
	})
}
