// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type Platform struct {
	Model

	Name          string `gorm:"unique"`
	RepoURL       string
	RepoDirectory string
	CommitSHA     string
	RefLabel      string // can be tag/branch/commit that user wrote in the yaml. example v0.1 or main.
	LabelType     int    // 1=branch, 2=tag, 3=commit

	Components []PlatformComponents `gorm:"foreignKey:PlatformID"` // Attributes of the module
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreatePlatform(p *Platform) (uuid.UUID, error) {
	return createOrUpdate(db.g(), p, []string{"name"})
}

func (pl *Platform) GetCondition() entity {
	return &Platform{
		Name: pl.Name,
	}
}
