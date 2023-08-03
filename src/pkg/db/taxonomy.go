package db

import "github.com/google/uuid"

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

func (t *Taxonomy) GetCondition() entity {
	return &Taxonomy{
		Level1: t.Level1,
		Level2: t.Level2,
		Level3: t.Level3,
		Level4: t.Level4,
		Level5: t.Level5,
		Level6: t.Level6,
		Level7: t.Level7,
	}
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTaxonomy(e *Taxonomy) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"level1", "level2", "level3", "level4", "level5", "level6", "level7"})
}
