package terrariumsrv

import (
	"context"

	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
)

func (s Service) ListModules(ctx context.Context, req *terrariumpb.ListModulesRequest) (resp *terrariumpb.ListModulesResponse, err error) {
	req.Page = setDefaultPage(req.Page)

	result, count, err := s.db.ListTFModule(req.Search, int(req.Page.Size), int(req.Page.Index*req.Page.Size))
	if err != nil {
		return nil, err
	}

	return &terrariumpb.ListModulesResponse{
		Page:    setPageResp(req.Page, count),
		Modules: result.ToProto(),
	}, nil
}
