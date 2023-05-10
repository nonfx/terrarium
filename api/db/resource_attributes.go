package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TFResourceAttribute struct {
	gorm.Model

	ResourceTypeID uint           `gorm:"uniqueIndex:resource_attribute_unique"`
	ResourceType   TFResourceType `gorm:"foreignKey:ResourceTypeID"`
	ProviderID     uint           `gorm:"uniqueIndex:resource_attribute_unique"`
	Provider       TFProvider     `gorm:"foreignKey:ProviderID"`
	AttributePath  string         `gorm:"uniqueIndex:resource_attribute_unique"`
	DataType       string
	Description    string
	Optional       bool
	Computed       bool
}

func (e *TFResourceAttribute) Create(g *gorm.DB) error {
	return g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "attribute_path"}, {Name: "provider_id"}, {Name: "resource_type_id"}},
		UpdateAll: true,
	}).Create(e).Error
}
