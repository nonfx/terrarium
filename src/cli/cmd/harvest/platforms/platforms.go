// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platforms

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/git"
	"gopkg.in/yaml.v2"

	terrpb "github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
)

func harvestPlatforms(g db.DB, directoryPath string) error {
	platforms, err := parsePlatformYAML(directoryPath)
	if err != nil {
		return err
	}

	for _, platform := range platforms {
		owner, repo, _, err := getOwnerRepoRef(platform)
		if err != nil {
			return err
		}

		for _, revision := range platform.Revisions {
			commitSHA, err := fetchCommitSHA(owner, repo, revision.Label)
			if err != nil {
				return err
			}

			dbPlatform := db.Platform{
				Title:         platform.Title,
				Description:   platform.Description,
				RepoURL:       platform.RepoURL,
				RepoDirectory: platform.RepoDir,
				CommitSHA:     commitSHA,
				RefLabel:      revision.Label,
				LabelType:     terrpb.GitLabelEnum(GitLabelEnumFromYAML(revision.Type)),
			}

			_, err = g.CreatePlatform(&dbPlatform)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parsePlatformYAML(directoryPath string) ([]Platform, error) {
	yamlContents, err := readPlatformYAML(directoryPath)
	if err != nil {
		return nil, err
	}

	var config struct {
		Platforms []Platform `yaml:"platforms"`
	}

	if err := yaml.Unmarshal([]byte(yamlContents), &config); err != nil {
		return nil, err
	}

	return config.Platforms, nil
}

type Platform struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	RepoURL     string     `yaml:"repo_url"`
	RepoDir     string     `yaml:"repo_directory"`
	Revisions   []Revision `yaml:"revisions"`
}

type Revision struct {
	Label string `yaml:"label"`
	Type  string `yaml:"type"`
}

func getOwnerRepoRef(platform Platform) (owner, repo, ref string, err error) {
	// Parse owner and repo from RepoURL
	parts := strings.Split(platform.RepoURL, "/")
	if len(parts) < 4 {
		return "", "", "", errors.New("invalid RepoURL format")
	}
	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]

	// Ref is the label
	ref = platform.Revisions[0].Label

	return owner, repo, ref, nil
}

func fetchCommitSHA(owner, repo, ref string) (string, error) {
	token := config.GitPassword()

	commitFetcher := git.NewGitHubCommitFetcher(token)

	sha, err := commitFetcher.FetchCommitSHA(owner, repo, ref)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	fmt.Println("Commit SHA:", ref, sha)
	return sha, nil
}

func readPlatformYAML(directoryPath string) (string, error) {
	fileInfo, err := os.Stat(directoryPath)
	if err != nil {
		return "", err
	}

	if fileInfo.IsDir() {
		yamlPath, err := findPlatformYAML(directoryPath)
		if err != nil {
			return "", err
		}

		yamlContents, err := os.ReadFile(yamlPath)
		if err != nil {
			return "", err
		}

		return string(yamlContents), nil
	} else if isPlatformYAML(fileInfo.Name()) {
		yamlContents, err := os.ReadFile(directoryPath)
		if err != nil {
			return "", err
		}
		fmt.Println(string(yamlContents))
		return string(yamlContents), nil
	} else {
		return "", fmt.Errorf("invalid file or directory: %s", directoryPath)
	}
}

func isPlatformYAML(filename string) bool {
	return filename == "platform.yaml" || filename == "platform.yml"
}

func findPlatformYAML(directoryPath string) (string, error) {
	var foundPath string

	extensions := map[string]bool{".yaml": true, ".yml": true}

	err := filepath.WalkDir(directoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if _, ok := extensions[ext]; ok {
			foundPath = path
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("platform.yaml or platform.yml not found in directory: %s", directoryPath)
	}

	return foundPath, nil
}

func GitLabelEnumFromYAML(language string) int32 {
	switch strings.ToLower(language) {
	case "branch":
		return int32(terrpb.GitLabelEnum_label_branch)
	case "tag":
		return int32(terrpb.GitLabelEnum_label_tag)
	case "commit":
		return int32(terrpb.GitLabelEnum_label_commit)
	default:
		return int32(terrpb.GitLabelEnum_label_no)
	}
}
