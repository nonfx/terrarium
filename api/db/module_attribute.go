package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TFModuleAttribute struct {
	gorm.Model

	ModuleID                       uint
	Module                         TFModule `gorm:"foreignKey:ModuleID"`
	ModuleAttributeName            string   `gorm:"uniqueIndex:module_attribute_unique"`
	Description                    string
	RelatedResourceTypeAttributeID uint
	ResourceAttribute              TFResourceAttribute `gorm:"foreignKey:RelatedResourceTypeAttributeID"`
	Optional                       bool
	Computed                       bool
}

func (e *TFModuleAttribute) Create(g *gorm.DB) error {
	return g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_attribute_name"}},
		UpdateAll: true,
	}).Create(e).Error
}
