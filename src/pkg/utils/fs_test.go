// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

var homeDir string = func() string {
	d, _ := homedir.Dir()
	return d
}()

var curDir string = func() string {
	d, _ := filepath.Abs(".")
	return d
}()

func TestResolveHomeAbs(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"~/.terrarium/test", homeDir + "/.terrarium/test"},
		{"./.terrarium/test", curDir + "/.terrarium/test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual, err := ResolveHomeAbs(tt.input)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSetupDir(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"~/.terrarium/test", homeDir + "/.terrarium/test"},
		{"./.terrarium/test", curDir + "/.terrarium/test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual, err := SetupDir(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)

			fi, err := os.Stat(actual)
			assert.NoError(t, err)
			assert.True(t, fi.IsDir())
		})
	}
}
