package db

import "gorm.io/gorm"

type TFResourceAttributesMapping struct {
	gorm.Model

	InputAttributeID  uint
	InputAttribute    TFResourceAttribute `gorm:"foreignKey:InputAttributeID"`
	OutputAttributeID uint
	OutputAttribute   TFResourceAttribute `gorm:"foreignKey:OutputAttributeID"`
}
