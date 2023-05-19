package db

import (
	"github.com/google/uuid"
)

type TFModuleAttribute struct {
	Model

	ModuleID                       uuid.UUID `gorm:"uniqueIndex:module_attribute_unique"`
	ModuleAttributeName            string    `gorm:"uniqueIndex:module_attribute_unique"`
	Description                    string
	RelatedResourceTypeAttributeID uuid.UUID
	Optional                       bool
	Computed                       bool

	Module            TFModule             `gorm:"foreignKey:ModuleID"`
	ResourceAttribute *TFResourceAttribute `gorm:"foreignKey:RelatedResourceTypeAttributeID"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFModuleAttribute(e *TFModuleAttribute) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"module_id", "module_attribute_name"})
}
