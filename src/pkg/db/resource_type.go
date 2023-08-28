package db

import "github.com/google/uuid"

type TFResourceType struct {
	Model

	ProviderID   uuid.UUID `gorm:"uniqueIndex:resource_type_unique"`
	ResourceType string    `gorm:"uniqueIndex:resource_type_unique"`
	TaxonomyID   uuid.UUID `gorm:"default:null"`

	Provider TFProvider `gorm:"foreignKey:ProviderID"`
	Taxonomy *Taxonomy  `gorm:"foreignKey:TaxonomyID"`
}

func (rt *TFResourceType) GetCondition() entity {
	return &TFResourceType{
		ProviderID:   rt.ProviderID,
		ResourceType: rt.ResourceType,
	}
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFResourceType(e *TFResourceType) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"provider_id", "resource_type"})
}

func (db *gDB) GetTFResourceType(e *TFResourceType, where *TFResourceType) error {
	return get(db.g(), e, where)
}
