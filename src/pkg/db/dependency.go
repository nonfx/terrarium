package db

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/google/uuid"
)

type Dependency struct {
	Model

	TaxonomyID  uuid.UUID        `gorm:"default:null"` // Given taxonomy's uncertain presence in YAML, setting TaxonomyID default as NULL accommodates potential absence of taxonomy data.
	InterfaceID string           `gorm:"unique"`
	Title       string           `gorm:"default:null"`
	Description string           `gorm:"default:null"`
	Inputs      *jsonschema.Node `gorm:"type:json"`
	Outputs     *jsonschema.Node `gorm:"type:json"`
	ExtendsID   string           `gorm:"-"` //This is yet to be finalized

	Taxonomy *Taxonomy `gorm:"foreignKey:TaxonomyID"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"interface_id"})
}

func (dep *Dependency) GetCondition() entity {
	return &Dependency{
		InterfaceID: dep.InterfaceID,
	}
}
