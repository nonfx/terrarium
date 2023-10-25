// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var re = regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)

type UpdateInfo struct {
	UpdateRequired       bool
	CurrentVersion       string
	LatestReleaseVersion string
}

func IsFarmUpdateRequired(cmd *cobra.Command, args []string) {
	updateInfo, _ := isUpdateRequired(config.FarmDefault())
	if updateInfo != nil && updateInfo.UpdateRequired {
		log.Info("farm update available", "currentVersion", updateInfo.CurrentVersion, "latestVersion", updateInfo.LatestReleaseVersion)
	}
}

// SetCurrentVersion sets the current version of farm repo in DB
func SetCurrentFarmVersion(version *db.FarmRelease) {
	g, err := config.DBConnect()
	if err != nil {
		log.Error("failed to connect DB", "error", err)
		return
	}
	g.CreateRelease(version)
}

func isUpdateRequired(repo string) (*UpdateInfo, error) {
	currentReleaseTag, err := GetCurrentReleaseTag(repo)
	if err != nil {
		return nil, err
	}

	latestReleaseTag, err := GetLatestReleaseTag(repo)
	if err != nil {
		return nil, err
	}

	if currentReleaseTag != latestReleaseTag {
		return &UpdateInfo{
			UpdateRequired:       true,
			CurrentVersion:       currentReleaseTag,
			LatestReleaseVersion: latestReleaseTag,
		}, nil
	}
	return &UpdateInfo{
		UpdateRequired: false,
	}, nil

}

func GetCurrentReleaseTag(repo string) (string, error) {
	g, err := config.DBConnect()
	if err != nil {
		return "", eris.Wrap(err, "failed to connect DB")
	}
	res := &db.FarmRelease{}
	err = g.FindReleaseByRepo(res, repo)
	if err != nil {
		return "", err
	}
	return res.Tag, nil
}

type Release struct {
	TagName string `json:"tag_name"`
}

func GetLatestReleaseTag(repoURL string) (string, error) {
	owner, repo, err := getOwnerAndRepo(repoURL)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", err
	}

	return release.TagName, nil
}

func getOwnerAndRepo(repoURL string) (owner, repo string, err error) {
	repoURL = strings.TrimSuffix(repoURL, ".git")
	if re.Match([]byte(repoURL)) {
		owner = string(re.ReplaceAll([]byte(repoURL), []byte("$1")))
		repo = string(re.ReplaceAll([]byte(repoURL), []byte("$2")))
	} else {
		err = fmt.Errorf("invalid repo")
	}
	return
}
