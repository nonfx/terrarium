package terrariumsrv

import (
	"context"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidRequest = status.Error(codes.InvalidArgument, "mappings cannot be populated with page size > 10")
)

func (s Service) ListModules(ctx context.Context, req *terrariumpb.ListModulesRequest) (resp *terrariumpb.ListModulesResponse, err error) {
	req.Page = setDefaultPage(req.Page)

	if req.Page.Size > 10 && req.PopulateMappings {
		return nil, eris.Wrapf(ErrInvalidRequest, "got page size: %d", req.Page.Size)
	}

	result, err := s.db.QueryTFModules(
		db.ModuleSearchFilter(req.Search),
		db.PopulateModuleMappingsFilter(req.PopulateMappings),
		db.PaginateGlobalFilter(req.Page.Size, req.Page.Index, &req.Page.Total),
	)
	if err != nil {
		return nil, err
	}

	return &terrariumpb.ListModulesResponse{
		Page:    req.Page,
		Modules: result.ToProto(),
	}, nil
}
