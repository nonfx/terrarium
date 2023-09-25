// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

type GitRef struct {
	Branch *string `json:"branch"`
	Tag    *string `json:"tag"`
	Commit *string `json:"commit"`
}

type GitVersion struct {
	Number string `json:"number"`
	GitRef GitRef `json:"git_ref"`
}

type Repository struct {
	URL         string       `json:"url"`
	Directory   string       `json:"directory"`
	GitVersions []GitVersion `json:"gitversions"`
}

type Platform struct {
	Model

	Name         string `gorm:"column:name"`
	Repositories string `gorm:"column:repositories;type:TEXT"` // JSON string of repositories
}
