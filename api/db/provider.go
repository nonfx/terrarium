package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TFProvider struct {
	gorm.Model

	Name string `gorm:"unique"`
}

func (e *TFProvider) Create(g *gorm.DB) error {
	return g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(e).Error
}
