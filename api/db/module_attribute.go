package db

import "gorm.io/gorm"

type TFModuleAttribute struct {
	gorm.Model

	ModuleID                       uint
	Module                         TFModule `gorm:"foreignKey:ModuleID"`
	ModuleAttributeName            string
	Description                    string
	RelatedResourceTypeAttributeID uint
	ResourceAttribute              TFResourceAttribute `gorm:"foreignKey:RelatedResourceTypeAttributeID"`
	Optional                       bool
	Computed                       bool
}
