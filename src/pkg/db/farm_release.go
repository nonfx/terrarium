// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type FarmRelease struct {
	Model
	Tag  string
	Repo string
}

func (fr *FarmRelease) GetCondition() entity {
	return &FarmRelease{
		Repo: fr.Repo,
	}
}

func (db *gDB) CreateRelease(e *FarmRelease) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"repo"})
}

func (db *gDB) FindReleaseByRepo(e *FarmRelease, repo string) error {
	return get(db.g(), e, &FarmRelease{Repo: repo})
}
