// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package git

import "github.com/google/go-github/github"

//go:generate mockery --name Git
type Git interface {
	FetchCommitSHA(owner, repo, ref string) (string, error)
	GetContents(owner, repo, ref, path string) (*github.RepositoryContent, []*github.RepositoryContent, error)
}
