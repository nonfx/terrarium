// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/google/uuid"
)

type TFProvider struct {
	Model

	Name string `gorm:"unique"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFProvider(e *TFProvider) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), e, []string{"name"})
	return id, err
}

func (db *gDB) GetTFProvider(e *TFProvider, where *TFProvider) error {
	return db.g().First(e, where).Error
}

func (db *gDB) GetOrCreateTFProvider(e *TFProvider) (id uuid.UUID, isNew bool, err error) {
	id, isNew, _, err = createOrGetOrUpdate(db.g(), e, []string{"name"})
	return
}

func (p1 *TFProvider) IsEq(p2 *TFProvider) bool {
	return p1.Name == p2.Name
}
