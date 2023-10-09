// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitClient struct {
	client *github.Client
}

func GithubClient(token string) Git {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GitClient{client: client}
}

func (gf *GitClient) FetchCommitSHA(owner, repo, ref string) (string, error) {
	commit, _, err := gf.client.Repositories.GetCommitSHA1(context.Background(), owner, repo, ref, "")
	if err != nil {
		return "", err
	}
	return commit, nil
}

func (gf *GitClient) GetContents(owner, repo, ref, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	content, list, _, err := gf.client.Repositories.GetContents(context.Background(), owner, repo, path, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return nil, nil, err
	}
	return content, list, nil
}
