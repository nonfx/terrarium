package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TFResourceType struct {
	gorm.Model

	ProviderID   uint       `gorm:"uniqueIndex:resource_type_unique"`
	Provider     TFProvider `gorm:"foreignKey:ProviderID"`
	ResourceType string     `gorm:"uniqueIndex:resource_type_unique"`
	TaxonomyID   string
}

func (e *TFResourceType) Create(g *gorm.DB) error {
	return g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_type"}, {Name: "provider_id"}},
		UpdateAll: true,
	}).Create(e).Error
}
