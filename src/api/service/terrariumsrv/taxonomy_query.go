// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"context"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
)

func (s Service) ListTaxonomy(ctx context.Context, req *terrariumpb.ListTaxonomyRequest) (resp *terrariumpb.ListTaxonomyResponse, err error) {
	req.Page = setDefaultPage(req.Page)

	result, err := s.db.QueryTaxonomies(db.TaxonomyRequestToFilters(req)...)
	if err != nil {
		return nil, eris.Wrap(err, "error running database query")
	}

	return &terrariumpb.ListTaxonomyResponse{
		Page:     req.Page,
		Taxonomy: result.ToProto(),
	}, nil
}
