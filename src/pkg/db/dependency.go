package db

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/google/uuid"
)

type Dependency struct {
	Model

	TaxonomyID  uuid.UUID        `json:"-" gorm:"unique"`
	Title       string           `json:"title" gorm:"unique"`
	Description string           `json:"description"`
	Inputs      *jsonschema.Node `gorm:"type:json"`
	Outputs     *jsonschema.Node `gorm:"type:json"`
	ExtendsID   string           `json:"extends_id" gorm:"-"` //This is yet to be finalized

	Taxonomy *Taxonomy `gorm:"foreignKey:TaxonomyID"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"title"})
}
