// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func (s Service) ListModuleAttributes(ctx context.Context, req *terrariumpb.ListModuleAttributesRequest) (*terrariumpb.ListModuleAttributesResponse, error) {
	req.Page = setDefaultPage(req.Page)

	if req.Page.Size > 100 && req.PopulateMappings {
		return nil, eris.Wrapf(ErrInvalidRequest, "got page size: %d", req.Page.Size)
	}

	result, err := s.db.QueryTFModuleAttributes(
		db.ModuleAttrByIDsFilter(uuid.MustParse(req.ModuleId)),
		db.ModuleAttrSearchFilter(req.Search),
		db.PopulateModuleAttrMappingsFilter(req.PopulateMappings),
		db.PaginateGlobalFilter(req.Page.Size, req.Page.Index, &req.Page.Total),
	)
	if err != nil {
		return nil, err
	}

	return &terrariumpb.ListModuleAttributesResponse{
		Page:       req.Page,
		Attributes: result.ToProto(),
	}, nil
}
