// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type FarmRelease struct {
	Model
	Tag  string
	Repo string `gorm:"unique"`
}

func (db *gDB) CreateRelease(e *FarmRelease) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), e, []string{"repo"})
	return id, err
}

func (db *gDB) FindReleaseByRepo(e *FarmRelease, repo string) error {
	return db.g().First(e, &FarmRelease{Repo: repo}).Error
}
