// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"gopkg.in/h2non/gock.v1"
)

func TestIsFarmUpdateRequired(t *testing.T) {
	config.LoadDefaults()
	mockDB := &mocks.DB{}
	mockDB.On("FindReleaseByRepo", mock.Anything, mock.Anything).Return(nil).Once()
	mockDB.On("FindReleaseByRepo", mock.Anything, mock.Anything).Return(func(e *db.FarmRelease, repo string) error {
		e.Tag = "mock_tag"
		return nil
	})
	config.SetDBMocks(mockDB)
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name      string
		args      args
		Gocksetup func()
	}{
		{
			name: "success with farm update required",
			args: args{
				cmd: &cobra.Command{},
			},
			Gocksetup: func() {
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
			},
		},
		{
			name: "success with farm update not required",
			args: args{
				cmd: &cobra.Command{},
			},
			Gocksetup: func() {
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
			},
		},
	}
	for _, tt := range tests {
		defer gock.Off()
		tt.Gocksetup()
		t.Run(tt.name, func(t *testing.T) {
			IsFarmUpdateRequired(tt.args.cmd, tt.args.args)
		})
	}
}

func TestSetCurrentFarmVersion(t *testing.T) {
	config.LoadDefaults()
	mockDB := &mocks.DB{}
	mockDB.On("CreateRelease", mock.Anything).Return(nil, nil)
	config.SetDBMocks(mockDB)
	type args struct {
		version *db.FarmRelease
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				version: &db.FarmRelease{},
			},
		},
		{
			name: "failed to connect to DB",
			args: args{
				version: &db.FarmRelease{},
			},
		},
	}
	for _, tt := range tests {
		if tt.name == "failed to connect to DB" {
			config.SetDBMocks(nil)
		}
		t.Run(tt.name, func(t *testing.T) {
			SetCurrentFarmVersion(tt.args.version)
		})
	}
}

func TestGetLatestReleaseTag(t *testing.T) {
	type args struct {
		repoURL string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantErr   bool
		Gocksetup func()
	}{
		{
			name: "invalid repo URL",
			args: args{
				repoURL: "invalid repo URL",
			},
			wantErr: true,
		},
		{
			name: "failed while invoking github API",
			args: args{
				repoURL: "invalid-url//github.com/cldcvr/terrarium-farm",
			},
			Gocksetup: func() {
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
			},
			wantErr: true,
		},
		{
			name: "non ok response from github API",
			args: args{
				repoURL: "github.com/cldcvr/terrarium-farm",
			},
			Gocksetup: func() {
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusInternalServerError).
					Body(bytes.NewReader(b))
			},
			wantErr: true,
		},
		{
			name: "failed to unmarshal response",
			args: args{
				repoURL: "github.com/cldcvr/terrarium-farm",
			},
			Gocksetup: func() {
				b, _ := json.Marshal("mock text")
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				repoURL: "github.com/cldcvr/terrarium-farm",
			},
			Gocksetup: func() {
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
			},
			want: "mock_tag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.Gocksetup != nil {
				defer gock.Off()
				tt.Gocksetup()
			}
			got, err := GetLatestReleaseTag(tt.args.repoURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestReleaseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLatestReleaseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentReleaseTag(t *testing.T) {
	config.LoadDefaults()

	type args struct {
		repo string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "failed to fetch release from DB",
			wantErr: true,
		},
		{
			name: "success",
			want: "mock_tag",
		},
		{
			name:    "failed to connect to DB",
			wantErr: true,
		},
	}
	mockDB := &mocks.DB{}
	mockDB.On("FindReleaseByRepo", mock.Anything, mock.Anything).Return(fmt.Errorf("mock error")).Once()
	mockDB.On("FindReleaseByRepo", mock.Anything, mock.Anything).Return(func(e *db.FarmRelease, repo string) error {
		e.Tag = "mock_tag"
		return nil
	})
	config.SetDBMocks(mockDB)
	for _, tt := range tests {
		if tt.name == "failed to connect to DB" {
			config.SetDBMocks(nil)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrentReleaseTag(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentReleaseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentReleaseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
