// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platforms

import (
	"context"
	"encoding/base64"
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
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	terrpb "github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
)

const (
	PlatformYAML = "platform.yaml"
	PlatformYML  = "platform.yml"
)

func harvestPlatforms(g db.DB, directoryPath string) error {
	platforms, err := parsePlatformYAML(directoryPath)
	if err != nil {
		return eris.Wrap(err, "error parsing platform YAML")
	}

	for _, platform := range platforms {
		err := processPlatform(g, platform)
		if err != nil {
			return err
		}
	}

	return nil
}

func processPlatform(g db.DB, platform Platform) error {
	owner, repo, _, err := getOwnerRepoRef(platform)
	if err != nil {
		return err
	}

	for _, revision := range platform.Revisions {
		err = processRevision(g, platform, owner, repo, revision)
		if err != nil {
			return err
		}
	}

	return nil
}

func processRevision(g db.DB, platform Platform, owner, repo string, revision Revision) error {
	ctx := context.Background()

	commitSHA, err := fetchCommitSHA(ctx, owner, repo, revision.Label)
	if err != nil {
		return err
	}

	dbPlatform := createDBPlatform(platform, commitSHA, revision)

	if _, err := g.CreatePlatform(&dbPlatform); err != nil {
		return err
	}

	terrariumYAMLPath, err := findTerrariumYAMLInGitHubDir(ctx, owner, repo, revision.Label, platform.RepoDir)
	if err != nil {
		return err
	}

	gc, _, err := gitGetContents(ctx, owner, repo, revision.Label, terrariumYAMLPath)
	if err != nil {
		return err
	}

	decodedContent, err := base64.StdEncoding.DecodeString(*gc.Content)
	if err != nil {
		return err
	}

	data, err := parseYAML(decodedContent)
	if err != nil {
		return err
	}

	// Fetch all the dependency id and interface id from the table
	q := g.Fetchdeps()
	for _, c := range data.Components {
		compareYAMLWithSQLData(g, c, q, dbPlatform.ID)
	}

	return nil
}

func createDBPlatform(platform Platform, commitSHA string, revision Revision) db.Platform {
	return db.Platform{
		Title:         platform.Title,
		Description:   platform.Description,
		RepoURL:       platform.RepoURL,
		RepoDirectory: platform.RepoDir,
		CommitSHA:     commitSHA,
		RefLabel:      revision.Label,
		LabelType:     terrpb.GitLabelEnum(GitLabelEnumFromYAML(revision.Type)),
	}
}

type YAMLData struct {
	Components []Component `yaml:"components"`
}

func parseYAML(decodedContent []byte) (YAMLData, error) {
	var data YAMLData

	if err := yaml.Unmarshal([]byte(decodedContent), &data); err != nil {
		return data, err
	}

	return data, nil
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

func gitClient(ctx context.Context) git.Git {
	t := config.GitToken()
	return git.GithubClient(ctx, t)
}

func fetchCommitSHA(ctx context.Context, owner, repo, ref string) (string, error) {
	sha, err := gitClient(ctx).FetchCommitSHA(ctx, owner, repo, ref)
	if err != nil {
		return "", err
	}
	return sha, nil
}

func gitGetContents(ctx context.Context, owner, repo, ref, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	gc, gl, err := gitClient(ctx).GetContents(ctx, owner, repo, ref, path)
	if err != nil {
		return nil, nil, err
	}
	return gc, gl, nil
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
		return string(yamlContents), nil
	} else {
		return "", fmt.Errorf("invalid file or directory: %s", directoryPath)
	}
}

func isPlatformYAML(filename string) bool {
	return filename == PlatformYAML || filename == PlatformYML
}

func findPlatformYAML(directoryPath string) (string, error) {
	var foundPath string

	filenames := []string{PlatformYAML, PlatformYML}

	err := filepath.WalkDir(directoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if slices.Contains(filenames, d.Name()) {
			foundPath = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("neither 'platform.yaml' nor 'platform.yml' was found in directory: %s", directoryPath)
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

func findTerrariumYAMLInGitHubDir(ctx context.Context, owner, repo, reference, dirPath string) (string, error) {
	gc, gl, err := gitGetContents(ctx, owner, repo, reference, dirPath)
	if err != nil {
		return "", err
	}

	if gc == nil {
		p, err := findTerrariumYaml(ctx, gl, owner, repo, reference, dirPath)
		if err != nil {
			return "", err
		}
		return p, nil
	}

	return "", fmt.Errorf("terrarium.yaml or terrarium.yml is not found in %s in directory: %s", reference, dirPath)
}

func findTerrariumYaml(ctx context.Context, gl []*github.RepositoryContent, owner, repo, reference, dirPath string) (string, error) {
	for _, content := range gl {
		if content.Type == nil {
			continue
		}

		if *content.Type == "file" && (strings.EqualFold(*content.Name, "terrarium.yaml") || strings.EqualFold(*content.Name, "terrarium.yml")) {
			return *content.Path, nil
		}

		if *content.Type == "dir" {
			subdirPath := *content.Path
			subfilePath, err := findTerrariumYAMLInGitHubDir(ctx, owner, repo, reference, subdirPath)
			if err != nil {
				return "", nil
			}
			return subfilePath, err
		}
	}
	return "", fmt.Errorf("terrarium.yaml or terrarium.yml is not found in %s in directory: %s", reference, dirPath)
}

func compareYAMLWithSQLData(g db.DB, c Component, queryResult []db.DependencyResult, u uuid.UUID) error {
	for _, r := range queryResult {
		if c.ID == r.InterfaceID {
			_, err := g.CreatePlatformComponents(&db.PlatformComponent{
				PlatformID:   u,
				DependencyID: r.DependencyID,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
