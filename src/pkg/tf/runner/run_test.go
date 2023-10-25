// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package runner

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/commander"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockExec struct {
	asserts    func(*exec.Cmd)
	assertsErr func(*exec.Cmd) error
}

func (m *mockExec) Run(cmd *exec.Cmd) error {
	if m.assertsErr != nil {
		return m.assertsErr(cmd)
	}
	if m.asserts != nil {
		m.asserts(cmd)
		return nil
	}
	return nil
}

func Test_terraformRunner_RunTerraformVersion(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "version",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commander.SetCommander(&mockExec{
				asserts: func(cmd *exec.Cmd) {
					assert.True(t, strings.HasSuffix(cmd.String(), "terraform version"))
					assert.Same(t, os.Stdout, cmd.Stdout)
				},
			})
			tr := NewTerraformRunner()
			err := tr.RunTerraformVersion()
			assert.NoError(t, err)
		})
	}
}

func Test_terraformRunner_RunTerraformInit(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "init",
			args: args{
				dir: "bypassing",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commander.SetCommander(&mockExec{
				asserts: func(cmd *exec.Cmd) {
					assert.True(t, strings.HasSuffix(cmd.String(), "terraform init"))
					assert.Same(t, os.Stdout, cmd.Stdout)
					assert.Equal(t, tt.args.dir, cmd.Dir)
				},
			})
			tr := NewTerraformRunner()
			err := tr.RunTerraformInit(tt.args.dir)
			assert.NoError(t, err)
		})
	}
}

func Test_terraformRunner_RunTerraformProviders(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "providers",
			args: args{
				dir: "Cambridgeshire",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commander.SetCommander(&mockExec{
				asserts: func(cmd *exec.Cmd) {
					assert.True(t, strings.HasSuffix(cmd.String(), "terraform providers"))
					assert.Same(t, os.Stdout, cmd.Stdout)
					assert.Equal(t, tt.args.dir, cmd.Dir)
				},
			})
			tr := NewTerraformRunner()
			err := tr.RunTerraformProviders(tt.args.dir)
			assert.NoError(t, err)
		})
	}
}

func Test_terraformRunner_RunTerraformProvidersSchema(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "*")
	dirInTmp := path.Join(tmpDir, "a_dir")
	fileInTmp := path.Join(dirInTmp, "a_file.json")
	require.NoError(t, err)
	require.NoError(t, os.Mkdir(dirInTmp, os.ModePerm))
	_, err = os.Create(fileInTmp)
	require.NoError(t, err)

	type args struct {
		dir         string
		outFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fail mkdir",
			args: args{
				outFilePath: path.Join(fileInTmp, "file_in_file.file"),
			},
			wantErr: true,
		},
		{
			name: "bad out-file",
			args: args{
				dir:         "Berkshire",
				outFilePath: dirInTmp, // dir instead of file
			},
			wantErr: true,
		},
		{
			name: "command error",
			args: args{
				dir:         "want-cmd-error",
				outFilePath: fileInTmp,
			},
			wantErr: true,
		},
		{
			name: "providers schema",
			args: args{
				dir:         "Lanka",
				outFilePath: fileInTmp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commander.SetCommander(&mockExec{
				assertsErr: func(cmd *exec.Cmd) error {
					if tt.args.dir == "want-cmd-error" {
						return eris.New("mocked error")
					}

					assert.True(t, strings.HasSuffix(cmd.String(), "terraform providers schema -json"))
					assert.Equal(t, tt.args.outFilePath, cmd.Stdout.(*os.File).Name())
					assert.Equal(t, tt.args.dir, cmd.Dir)
					return nil
				},
			})
			tr := NewTerraformRunner()
			err := tr.RunTerraformProvidersSchema(tt.args.dir, tt.args.outFilePath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
