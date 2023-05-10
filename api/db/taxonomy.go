package db

import "gorm.io/gorm"

type Taxonomy struct {
	gorm.Model

	TaxonomyLevel1 string
	TaxonomyLevel2 string
	TaxonomyLevel3 string
}
