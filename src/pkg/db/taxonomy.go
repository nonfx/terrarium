// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/google/uuid"
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
	id, _, _, err := createOrGetOrUpdate(db.g(), e, []string{"level1", "level2", "level3", "level4", "level5", "level6", "level7"})
	return id, err
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
