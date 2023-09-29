// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platforms

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/git"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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

			terrariumYAMLPath, err := findTerrariumYAMLInGitHubDir(owner, repo, revision.Label, platform.RepoDir, "")
			if err != nil {
				return err
			}

			fmt.Println(terrariumYAMLPath)

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

			// gc, err := gitGetContents(owner, repo, revision.Label, platform.RepoDir)
			// if err != nil {
			// 	return err
			// }
			// decodedContent, err := base64.StdEncoding.DecodeString(gc)
			// if err != nil {
			// 	fmt.Println("Error decoding content:", err)
			// 	return err
			// }
			// fmt.Println(string(decodedContent))

			// _, err = g.CreatePlatformComponents(&db.PlatformComponents{
			// 	PlatformID:   dbPlatform.ID,
			// 	DependencyID: db.Dependency{InterfaceID: "test"}.ID,
			// })
			// if err != nil {
			// 	return err
			// }
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

func gitClient() git.Git {
	t := config.GitPassword()
	return git.GithubClient(t)
}

func fetchCommitSHA(owner, repo, ref string) (string, error) {
	sha, err := gitClient().FetchCommitSHA(owner, repo, ref)
	if err != nil {
		return "", err
	}
	return sha, nil
}

func gitGetContents(owner, repo, ref, path string) (string, error) {
	gc, err := gitClient().GetContents(owner, repo, ref, path)
	if err != nil {
		return "", err
	}
	return gc, nil
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

// func findTerrariumYAML(owner, repo, ref, dir string) (string, error) {
// 	// Fetch the contents of the directory in the GitHub repository.
// 	contents, err := gitClient().GetContents(owner, repo, ref, dir)
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println(contents)

// yamlFilenames := []string{"terrarium.yaml", "terrarium.yml"}

// Iterate through the fetched contents to check for the YAML file.
// for _, content := range contents {
// 	if content.Type == "file" {
// 		for _, filename := range yamlFilenames {
// 			if content.Name == filename {
// 				return content.Path, nil
// 			}
// 		}
// 	}
// }

// 	return "", fmt.Errorf("terrarium.yaml or terrarium.yml not found in directory: %s", dir)
// }

// func findTerrariumYAMLPath(directoryPath, owner, repo, reference string) (string, error) {
// 	// Check if the directory path ends with terrarium.yaml or terrarium.yml
// 	if isTerrariumYAMLFile(directoryPath) {
// 		return directoryPath, nil
// 	}

// 	// Try appending terrarium.yaml to the directory path and check if it exists
// 	terrariumYAMLPath := filepath.Join(directoryPath, "terrarium.yaml")
// 	if fileExists(terrariumYAMLPath) {
// 		return terrariumYAMLPath, nil
// 	}

// 	// Try appending terrarium.yml to the directory path and check if it exists
// 	terrariumYMLPath := filepath.Join(directoryPath, "terrarium.yml")
// 	if fileExists(terrariumYMLPath) {
// 		return terrariumYMLPath, nil
// 	}

// 	// If none of the options exist, return an error
// 	return "", fmt.Errorf("terrarium.yaml or terrarium.yml not found in directory: %s", directoryPath)
// }

// func isTerrariumYAMLFile(filePath string) bool {
// 	return filepath.Ext(filePath) == ".yaml" || filepath.Ext(filePath) == ".yml"
// }

// func fileExists(filePath string) bool {
// 	_, err := os.Stat(filePath)
// 	return err == nil
// }

func findTerrariumYAMLInGitHubDir(owner, repo, reference, dirPath, token string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the contents of the directory in the GitHub repository as a list of RepositoryContent items.
	_, _, _, err := client.Repositories.GetContents(ctx, owner, repo, dirPath, &github.RepositoryContentGetOptions{
		Ref: reference,
	})
	if err != nil {
		return "", err
	}

	// Loop through the contents in the directory.
	// for _, content := range contents {
	// 	if content.Type != nil {
	// 		if *content.Type == "file" {
	// 			// Check if the file is named "terrarium.yaml" or "terrarium.yml".
	// 			if strings.EqualFold(*content.Name, "terrarium.yaml") || strings.EqualFold(*content.Name, "terrarium.yml") {
	// 				// If found, return the path to the file.
	// 				return *content.Path, nil
	// 			}
	// 		} else if *content.Type == "dir" {
	// 			// If it's a subdirectory, recursively search for the file within it.
	// 			subdirPath := *content.Path
	// 			subfilePath, err := findTerrariumYAMLInGitHubDir(owner, repo, reference, subdirPath, token)
	// 			if err == nil {
	// 				return subfilePath, nil
	// 			}
	// 		}
	// 	}
	// }

	// If the file is not found in the directory or its subdirectories, return an error.
	return "", fmt.Errorf("terrarium.yaml or terrarium.yml not found in directory: %s", dirPath)
}
