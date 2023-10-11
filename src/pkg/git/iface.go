// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"context"

	"github.com/google/go-github/github"
)

//go:generate mockery --name Git
type Git interface {
	FetchCommitSHA(ctx context.Context, owner, repo, ref string) (string, error)
	GetContents(ctx context.Context, owner, repo, ref, path string) (*github.RepositoryContent, []*github.RepositoryContent, error)
}
