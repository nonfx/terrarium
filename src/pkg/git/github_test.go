// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitClient(t *testing.T) {
	client := GithubClient("")
	owner := "cldcvr"
	repo := "terrarium"
	ref := "main"
	path := "examples/platform/terrarium.yaml"

	t.Run("FetchCommitSHA", func(t *testing.T) {
		commit, err := client.FetchCommitSHA(owner, repo, ref)
		if err != nil {
			t.Error("Error fetching commit SHA: ", err)
			return
		}
		assert.NotEmpty(t, commit)
	})

	t.Run("GetContents", func(t *testing.T) {
		content, _, err := client.GetContents(owner, repo, ref, path)
		if err != nil {
			t.Error("Error getting contents: ", err)
			return
		}
		assert.NotNil(t, content)
		assert.NotEmpty(t, content)
	})

	t.Run("FetchCommitSHA (error case)", func(t *testing.T) {
		commit, err := client.FetchCommitSHA(owner, repo, "non-existent-ref")
		if err == nil {
			t.Error("Expected error when fetching commit SHA")
			return
		}
		assert.Empty(t, commit)
	})

	t.Run("GetContents (error case)", func(t *testing.T) {
		content, _, err := client.GetContents(owner, repo, ref, "non-existent-path")
		if err == nil {
			t.Error("Expected error when getting contents")
			return
		}
		assert.Nil(t, content)
	})
}
