// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
)

func (s Service) ListComponents(ctx context.Context, req *terrariumpb.ListComponentsRequest) (resp *terrariumpb.ListComponentsResponse, err error) {
	req.Page = setDefaultPage(req.Page)

	result, err := s.db.QueryPlatformComponents(db.ComponentRequestToFilters(req)...)
	if err != nil {
		return nil, eris.Wrap(err, "error running database query")
	}

	return &terrariumpb.ListComponentsResponse{
		Page:       req.Page,
		Components: result.ToProto(),
	}, nil
}
