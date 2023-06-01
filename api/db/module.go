package db

import (
	"github.com/google/uuid"
)

type TFModule struct {
	Model

	ModuleName  string
	Source      string `gorm:"uniqueIndex:module_unique"`
	Version     string `gorm:"uniqueIndex:module_unique"`
	Description string
	TaxonomyID  uuid.UUID `gorm:"default:null"`

	Taxonomy *Taxonomy `gorm:"foreignKey:TaxonomyID"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFModule(e *TFModule) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"source", "version"})
}
