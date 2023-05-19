package db

import (
	"github.com/google/uuid"
)

type TFResourceAttributesMapping struct {
	Model

	InputAttributeID  uuid.UUID
	InputAttribute    TFResourceAttribute `gorm:"foreignKey:InputAttributeID"`
	OutputAttributeID uuid.UUID
	OutputAttribute   TFResourceAttribute `gorm:"foreignKey:OutputAttributeID"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFResourceAttributesMapping(e *TFResourceAttributesMapping) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"input_attribute_id", "output_attribute_id"})
}
