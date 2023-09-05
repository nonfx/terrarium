// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package runner

import (
	"os"
	"path"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/commander"
	"github.com/stretchr/testify/assert"
)

func TestTerraformProviderSchema(t *testing.T) {
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
			name: "providers schema",
			args: args{
				dir:         os.TempDir(),
				outFilePath: path.Join(os.TempDir(), "166e4331-7d85-4c0c-8d8a-3587a8ef95b9"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		commander.SetCommander(&mockExec{})
		tr := NewTerraformRunner()
		err := TerraformProviderSchema(tr, tt.args.dir, tt.args.outFilePath)
		assert.NoError(t, err)
	}
}

func TestTerraformInit(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "providers schema",
			args: args{
				dir: "Avon",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		commander.SetCommander(&mockExec{})
		tr := NewTerraformRunner()
		err := TerraformInit(tr, tt.args.dir)
		assert.NoError(t, err)
	}
}
