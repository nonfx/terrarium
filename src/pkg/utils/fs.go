// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/rotisserie/eris"
)

// ResolveHomeAbs resolves home directory or relative path into absolute path
func ResolveHomeAbs(path string) (absPath string, err error) {
	dir, err := homedir.Expand(path)
	if err != nil {
		return "", eris.Wrapf(err, "failed to resolve home dir in given path: %s", path)
	}

	absPath, err = filepath.Abs(dir)
	if err != nil {
		return "", eris.Wrapf(err, "failed to resolve into absolute path: %s", path)
	}

	return
}

// SetupDir resolves home dir (`~`) or relative path and then creates directory if missing.
func SetupDir(dirPath string) (absDirPath string, err error) {
	absDirPath, err = ResolveHomeAbs(dirPath)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(absDirPath, os.ModePerm)
	if err != nil {
		return "", eris.Wrapf(err, "failed to create dir: %s", absDirPath)
	}

	return
}
