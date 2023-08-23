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
	"github.com/stretchr/testify/assert"
)

type mockExec struct {
	asserts func(*exec.Cmd)
}

func (m *mockExec) Run(cmd *exec.Cmd) error {
	if m.asserts != nil {
		m.asserts(cmd)
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
	mockOut, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
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
			name: "bad out-file",
			args: args{
				dir:         "Berkshire",
				outFilePath: path.Join(os.TempDir(), "ae66e499-9185-4639-a02a-79f870e1dce9", "d4285cd7-915f-4aa8-a193-9e6db7a033f0"),
			},
			wantErr: true,
		},
		{
			name: "providers schema",
			args: args{
				dir:         "Lanka",
				outFilePath: mockOut.Name(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commander.SetCommander(&mockExec{
				asserts: func(cmd *exec.Cmd) {
					assert.True(t, strings.HasSuffix(cmd.String(), "terraform providers schema -json"))
					assert.Equal(t, tt.args.outFilePath, cmd.Stdout.(*os.File).Name())
					assert.Equal(t, tt.args.dir, cmd.Dir)
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
