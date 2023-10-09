// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type Taxonomy struct {
	Model

	Level1 string `gorm:"uniqueIndex:taxonomy_unique"`
	Level2 string `gorm:"uniqueIndex:taxonomy_unique"`
	Level3 string `gorm:"uniqueIndex:taxonomy_unique"`
	Level4 string `gorm:"uniqueIndex:taxonomy_unique"`
	Level5 string `gorm:"uniqueIndex:taxonomy_unique"`
	Level6 string `gorm:"uniqueIndex:taxonomy_unique"`
	Level7 string `gorm:"uniqueIndex:taxonomy_unique"`
}

type Taxonomies []Taxonomy

var taxonomyLevelCols = []string{"level1", "level2", "level3", "level4", "level5", "level6", "level7"}

func (t1 *Taxonomy) IsEq(t2 *Taxonomy) bool {
	return t1.Level1 == t2.Level1 &&
		t1.Level2 == t2.Level2 &&
		t1.Level3 == t2.Level3 &&
		t1.Level4 == t2.Level4 &&
		t1.Level5 == t2.Level5 &&
		t1.Level6 == t2.Level6 &&
		t1.Level7 == t2.Level7
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTaxonomy(e *Taxonomy) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), e, taxonomyLevelCols)
	return id, err
}

// QueryTaxonomies based on the given filters
func (db *gDB) QueryTaxonomies(filterOps ...FilterOption) (result Taxonomies, err error) {
	q := db.g().Model(&Taxonomy{}).Order(strings.Join(taxonomyLevelCols, ", ")).Not(&Taxonomy{Model: Model{ID: uuid.Nil}}, "id")

	for _, filer := range filterOps {
		q = filer(q)
	}

	err = q.Find(&result).Error
	if err != nil {
		return nil, eris.Wrap(err, "query taxonomy")
	}

	return
}

func TaxonomyByLevelsFilter(t *Taxonomy) FilterOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where(t)
	}
}

func TaxonomyFromLevels(levels ...string) *Taxonomy {
	t := &Taxonomy{}

	if len(levels) > 0 {
		t.Level1 = levels[0]
	}

	if len(levels) > 1 {
		t.Level2 = levels[1]
	}

	if len(levels) > 2 {
		t.Level3 = levels[2]
	}

	if len(levels) > 3 {
		t.Level4 = levels[3]
	}

	if len(levels) > 4 {
		t.Level5 = levels[4]
	}

	if len(levels) > 5 {
		t.Level6 = levels[5]
	}

	if len(levels) > 6 {
		t.Level7 = levels[6]
	}

	return t
}

func (t *Taxonomy) ToLevels() []string {
	levels := []string{
		t.Level1,
		t.Level2,
		t.Level3,
		t.Level4,
		t.Level5,
		t.Level6,
		t.Level7,
	}

	return utils.TrimEmpty(levels)
}

func (tArr Taxonomies) ToProto() []*terrariumpb.Taxonomy {
	resp := make([]*terrariumpb.Taxonomy, len(tArr))

	for i, t := range tArr {
		resp[i] = t.ToProto()
	}

	return resp
}

func (t *Taxonomy) ToProto() *terrariumpb.Taxonomy {
	return &terrariumpb.Taxonomy{
		Id:     t.ID.String(),
		Levels: t.ToLevels(),
	}
}

func TaxonomyRequestToFilters(req *terrariumpb.ListTaxonomyRequest) []FilterOption {
	filters := []FilterOption{}

	if req.Page != nil {
		filters = append(filters, PaginateGlobalFilter(req.Page.Size, req.Page.Index, &req.Page.Total))
	}

	if len(req.Taxonomy) > 0 {
		filters = append(filters, TaxonomyByLevelsFilter(TaxonomyFromLevels(req.Taxonomy...)))
	}

	return filters
}
