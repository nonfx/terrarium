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

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTaxonomy(e *Taxonomy) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"level1", "level2", "level3", "level4", "level5", "level6", "level7"})
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"interface_id"})
}
