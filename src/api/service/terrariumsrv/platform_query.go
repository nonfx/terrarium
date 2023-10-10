// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
)

func (s Service) ListPlatforms(ctx context.Context, req *terrariumpb.ListPlatformsRequest) (resp *terrariumpb.ListPlatformsResponse, err error) {
	req.Page = setDefaultPage(req.Page)

	result, err := s.db.QueryPlatforms(db.PlatformRequestToFilters(req)...)
	if err != nil {
		return nil, eris.Wrap(err, "error running database query")
	}

	return &terrariumpb.ListPlatformsResponse{
		Page:      req.Page,
		Platforms: result.ToProto(),
	}, nil
}
