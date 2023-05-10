package db

import "gorm.io/gorm"

type TFModule struct {
	gorm.Model

	TaxonomyID  uint
	Taxonomy    Taxonomy `gorm:"foreignKey:TaxonomyID"`
	ModuleName  string
	Source      string
	Description string
}
