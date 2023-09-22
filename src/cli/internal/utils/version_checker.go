// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/rotisserie/eris"
)

// add 15 minutes caching for farm release verion to ensure
// that we are not querying for latest release frequently

var re = regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)

type UpdateInfo struct {
	UpdateRequired       bool
	CurrentVersion       string
	LatestReleaseVersion string
}

// SetCurrentVersion sets the current version of farm repo
func SetCurrentFarmVersion(version *db.FarmRelease) {
	g, err := config.DBConnect()
	if err != nil {
		log.Default().Println(err, "failed to connect DB")
		return
	}
	g.CreateRelease(version)
}

func IsUpdateRequired(repo string) (*UpdateInfo, error) {
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

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response into a Release struct
	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", err
	}

	return release.TagName, nil
}

func getOwnerAndRepo(repoURL string) (owner, repo string, err error) {
	repoURL = strings.TrimSuffix(repoURL, ".git")
	if re.Match([]byte(repo)) {
		owner = string(re.ReplaceAll([]byte(repoURL), []byte("$1")))
		repo = string(re.ReplaceAll([]byte(repoURL), []byte("$2")))
	} else {
		err = fmt.Errorf("invalid repo")
	}
	return
}
