package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TFModule struct {
	gorm.Model

	TaxonomyID  uint
	Taxonomy    Taxonomy `gorm:"-"`
	ModuleName  string
	Source      string `gorm:"uniqueIndex:module_unique"`
	Description string
}

func (e *TFModule) Create(g *gorm.DB) error {
	return g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source"}},
		UpdateAll: true,
	}).Create(e).Error
}
