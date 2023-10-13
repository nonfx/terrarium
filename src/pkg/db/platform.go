// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/cldcvr/terrarium/src/pkg/metadata/taxonomy"
	terrpb "github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type Platform struct {
	Model

	Title         string
	Description   string
	RepoURL       string
	RepoDirectory string
	CommitSHA     string              `gorm:"unique"`
	RefLabel      string              // can be tag/branch/commit that user wrote in the yaml. example v0.1 or main.
	LabelType     terrpb.GitLabelEnum // 1=branch, 2=tag, 3=commit

	Components []PlatformComponent `gorm:"foreignKey:PlatformID"`
}

type Platforms []Platform

// CreatePlatform insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreatePlatform(p *Platform) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), p, []string{"commit_sha"})
	return id, err
}

// QueryPlatforms query platforms table and return list of matching platforms
func (db *gDB) QueryPlatforms(filterOps ...FilterOption) (result Platforms, err error) {
	q := db.g().Model(&Platform{})

	for _, filer := range filterOps {
		q = filer(q)
	}

	q = q.Order("title").Preload("Components")

	err = q.Find(&result).Error
	if err != nil {
		return nil, eris.Wrap(err, "query platform")
	}

	return
}

// PlatformFilterBySearch perform search on name & repo columns
func PlatformFilterBySearch(query string) FilterOption {
	if query == "" {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		q := "%" + query + "%"
		return g.Where("title LIKE ? OR repo_url LIKE ?", q, q)
	}
}

func PlatformFilterByTaxonomy(tax *Taxonomy) FilterOption {
	return func(g *gorm.DB) *gorm.DB {
		taxQ := g.Session(&gorm.Session{NewDB: true}).Model(&Taxonomy{}).Where(tax).Select("id")
		depQ := g.Session(&gorm.Session{NewDB: true}).Model(&Dependency{}).Where("taxonomy_id IN (?)", taxQ).Select("id")
		compQ := g.Session(&gorm.Session{NewDB: true}).Model(&PlatformComponent{}).Where("dependency_id IN (?)", depQ).Select("platform_id")
		return g.Where("id IN (?)", compQ)
	}
}

func PlatformFilterByDependencyID(depIDs ...string) FilterOption {
	if len(utils.TrimEmpty(depIDs)) == 0 {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		compQ := g.Session(&gorm.Session{NewDB: true}).Model(&PlatformComponent{}).Where("dependency_id IN (?)", depIDs).Select("platform_id")
		return g.Where("id IN (?)", compQ)
	}
}

func (p Platform) ToProto() *terrpb.Platform {
	return &terrpb.Platform{
		Id:         p.ID.String(),
		Title:      p.Title,
		RepoUrl:    p.RepoURL,
		RepoDir:    p.RepoDirectory,
		RepoCommit: p.CommitSHA,
		RefLabel:   p.RefLabel,
		RefType:    p.LabelType,
		Components: int32(len(p.Components)),
	}
}

func (pArr Platforms) ToProto() []*terrpb.Platform {
	respArr := make([]*terrpb.Platform, len(pArr))

	for i, p := range pArr {
		respArr[i] = p.ToProto()
	}

	return respArr
}
func PlatformRequestToFilters(req *terrpb.ListPlatformsRequest) []FilterOption {
	filters := []FilterOption{}

	if req.Page != nil {
		filters = append(filters, PaginateGlobalFilter(req.Page.Size, req.Page.Index, &req.Page.Total))
	}

	tax := taxonomy.Taxon(req.Taxonomy).Split()
	if len(tax) > 0 {
		filters = append(filters, PlatformFilterByTaxonomy(TaxonomyFromLevels(tax...)))
	}

	filters = append(filters,
		PlatformFilterByDependencyID(req.InterfaceUuid...),
		PlatformFilterBySearch(req.Search),
	)

	return filters
}
