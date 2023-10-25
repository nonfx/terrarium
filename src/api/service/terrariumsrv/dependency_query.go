// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
)

func (s Service) ListDependencies(ctx context.Context, req *terrariumpb.ListDependenciesRequest) (resp *terrariumpb.ListDependenciesResponse, err error) {
	req.Page = setDefaultPage(req.Page)

	result, err := s.db.QueryDependencies(db.DependencyRequestToFilters(req)...)
	if err != nil {
		return nil, eris.Wrap(err, "error running database query")
	}

	return &terrariumpb.ListDependenciesResponse{
		Page:         req.Page,
		Dependencies: result.ToProto(),
	}, nil
}
