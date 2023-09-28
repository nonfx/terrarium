// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// type GIT interface {
// 	FetchCommitSHA(owner, repo, ref string) (string, error)
// }

type GitHubCommitFetcher struct {
	client *github.Client
}

func NewGitHubCommitFetcher(token string) Git {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GitHubCommitFetcher{client: client}
}

func (gf *GitHubCommitFetcher) FetchCommitSHA(owner, repo, ref string) (string, error) {
	commit, _, err := gf.client.Repositories.GetCommitSHA1(context.Background(), owner, repo, ref, "")
	if err != nil {
		return "", err
	}
	fmt.Println(commit)
	return commit, nil
}

func (gf *GitHubCommitFetcher) GetContents(owner, repo, ref string) (string, error) {
	content, _, _, err := gf.client.Repositories.GetContents(context.Background(), owner, repo, "examples/platform/terrarium.yaml", &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return "", err
	}
	fmt.Println((*content.Content))
	decodedContent, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		fmt.Println("Error decoding content:", err)
		return "", err
	}
	fmt.Println(string(decodedContent))

	return "", nil
}
