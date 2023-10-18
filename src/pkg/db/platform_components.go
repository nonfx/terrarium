// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/cldcvr/terrarium/src/pkg/metadata/taxonomy"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type PlatformComponent struct {
	Model

	PlatformID   uuid.UUID `gorm:"uniqueIndex:platform_components_unique"`
	DependencyID uuid.UUID `gorm:"uniqueIndex:platform_components_unique"`

	Platform   Platform   `gorm:"foreignKey:PlatformID"`
	Dependency Dependency `gorm:"foreignKey:DependencyID"`
}

type PlatformComponents []PlatformComponent

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreatePlatformComponents(p *PlatformComponent) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), p, []string{"platform_id", "dependency_id"})
	return id, err
}

// QueryPlatforms query platforms table and return list of matching platforms
func (db *gDB) QueryPlatformComponents(filterOps ...FilterOption) (result PlatformComponents, err error) {
	q := db.g().Model(&PlatformComponents{})

	for _, filer := range filterOps {
		q = filer(q)
	}

	q = q.Order("created_at").Preload("Dependency.Taxonomy").Preload("Dependency.Attributes")

	err = q.Find(&result).Error
	if err != nil {
		return nil, eris.Wrap(err, "query platform components")
	}

	return
}

func ComponentsFilterByTaxonomy(tax *Taxonomy) FilterOption {
	return func(g *gorm.DB) *gorm.DB {
		taxQ := g.Session(&gorm.Session{NewDB: true}).Model(&Taxonomy{}).Where(tax).Select("id")
		depQ := g.Session(&gorm.Session{NewDB: true}).Model(&Dependency{}).Where("taxonomy_id IN (?)", taxQ).Select("id")
		return g.Where("dependency_id IN (?)", depQ)
	}
}

func ComponentsFilterByDependencySearch(query string) FilterOption {
	if query == "" {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		q := "%" + query + "%"
		depQ := g.Session(&gorm.Session{NewDB: true}).Model(&Dependency{}).Where("interface_id LIKE ? OR title LIKE ?", q, q).Select("id")
		return g.Where("dependency_id IN (?)", depQ)
	}
}

func ComponentsFilterByPlatformID(ids ...string) FilterOption {
	ids = utils.TrimEmpty(ids)
	if len(ids) == 0 {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		return g.Where("platform_id IN (?)", ids)
	}
}

func (c PlatformComponent) ToProto() *terrariumpb.Component {
	depProto := c.Dependency.ToProto()

	return &terrariumpb.Component{
		Id:            c.ID.String(),
		InterfaceUuid: c.DependencyID.String(),
		InterfaceId:   depProto.InterfaceId,
		Title:         depProto.Title,
		Description:   depProto.Description,
		Taxonomy:      depProto.Taxonomy,
		Inputs:        depProto.Inputs,
		Outputs:       depProto.Outputs,
	}
}

func (c PlatformComponents) ToProto() []*terrariumpb.Component {
	resp := make([]*terrariumpb.Component, len(c))
	for i, c := range c {
		resp[i] = c.ToProto()
	}

	return resp
}

func ComponentRequestToFilters(req *terrariumpb.ListComponentsRequest) []FilterOption {
	filters := []FilterOption{}

	if req.Page != nil {
		filters = append(filters, PaginateGlobalFilter(req.Page.Size, req.Page.Index, &req.Page.Total))
	}

	if req.Taxonomy != "" {
		tax := TaxonomyFromLevels(taxonomy.Taxon(req.Taxonomy).Split()...)
		filters = append(filters, ComponentsFilterByTaxonomy(tax))
	}

	filters = append(
		filters,
		ComponentsFilterByPlatformID(req.PlatformId),
		ComponentsFilterByDependencySearch(req.Search),
	)

	return filters
}
