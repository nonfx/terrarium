// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package update

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/cli/internal/utils"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

const (
	artifactName = "cc_terrarium_data.sql.gz"
)

var (
	cmd *cobra.Command
)

func NewCmd() *cobra.Command {
	cmd = &cobra.Command{
		Use:   "update",
		Short: "farm update updates the databse with latest farm release",
		RunE:  cmdRunE,
	}
	return cmd
}

func cmdRunE(cmd *cobra.Command, args []string) error {
	dumpFilePath, latestReleaseTag, err := downloadArtifact(config.FarmDefault(), artifactName)
	if err != nil {
		return err
	}

	gzFile, err := os.Open(filepath.Join(dumpFilePath, artifactName))
	if err != nil {
		return eris.Wrap(err, "error opening file: %v")
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return eris.Wrap(err, "error creating gzip reader: %v")
	}
	defer gzReader.Close()

	dumpContent, err := io.ReadAll(gzReader)
	if err != nil {
		return eris.Wrap(err, "error reading content: %v")
	}

	err = seedDatabase(string(dumpContent))
	if err != nil {
		return eris.Wrap(err, "failed to seed database: %v")
	}

	// we store the latest artifact as current running artifact
	utils.SetCurrentFarmVersion(&db.FarmRelease{
		Repo: config.FarmDefault(),
		Tag:  latestReleaseTag,
	})
	return nil
}

// downloadArtifact downloads the latest artifact
func downloadArtifact(repo, artifactName string) (string, string, error) {
	releaseTag, err := utils.GetLatestReleaseTag(repo)
	if err != nil {
		return "", "", eris.Wrap(err, "failed to fetch latest release tag: %v")
	}

	resp, err := http.Get(fmt.Sprintf("https://%s/releases/download/%s/%s", repo, releaseTag, artifactName))
	if err != nil {
		return "", "", eris.Wrap(err, "error sending HTTP request: %v")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", eris.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	tempDir := filepath.Join(os.TempDir(), "farm-artifact")
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", "", eris.Wrap(err, "error creating temporary directory: %v")
	}

	outFile, err := os.Create(filepath.Join(tempDir, artifactName))
	if err != nil {
		return "", "", eris.Wrap(err, "error creating output file: %v")
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", "", eris.Wrap(err, "error copying artifact to output file: %v")
	}
	return tempDir, releaseTag, nil
}
