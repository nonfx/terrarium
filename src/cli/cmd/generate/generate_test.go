// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_copyFile(t *testing.T) {
	type args struct {
		srcPath string
		dstPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fail on source (not exists)",
			args: args{
				srcPath: "/ggsgsd/rthrtyytjurtkuykkuyyukruy/cambridgeshire.emf.fnc",
				dstPath: mustCreateFile(t, []byte("Sint non incidunt ipsam et.\n")),
			},
			wantErr: true,
		},
		{
			name: "fail on source (not file)",
			args: args{
				srcPath: os.TempDir(),
				dstPath: mustCreateFile(t, []byte("Sint non incidunt ipsam et.\n")),
			},
			wantErr: true,
		},
		{
			name: "fail on destination (not file)",
			args: args{
				srcPath: mustCreateFile(t, []byte("Sint non incidunt ipsam et.\n")),
				dstPath: os.TempDir(),
			},
			wantErr: true,
		},
		{
			name: "successful copy (dest overwrite)",
			args: args{
				srcPath: mustCreateFile(t, []byte("Totam facere iusto laborum eum possimus sint harum.\nRepudiandae possimus nulla autem dolorem fuga veniam.\n")),
				dstPath: mustCreateFile(t, []byte("Sint non incidunt ipsam et.\n")),
			},
			wantErr: false,
		},
		{
			name: "successful copy (dest not exists)",
			args: args{
				srcPath: mustCreateFile(t, []byte("Sint non incidunt ipsam et.\n")),
				dstPath: path.Join(os.TempDir(), "generate.tmo"),
			},
			wantErr: false,
		},
		{
			name: "successful empty copy",
			args: args{
				srcPath: mustCreateFile(t, nil),
				dstPath: mustCreateFile(t, []byte("Sint non incidunt ipsam et.\n")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := copyFile(tt.args.srcPath, tt.args.dstPath); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assertFileSame(t, tt.args.srcPath, tt.args.dstPath)
			}
		})
	}
}

func mustCreateFile(t *testing.T, contents []byte) string {
	fp, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	if len(contents) > 0 {
		if _, err := fp.Write(contents); err != nil {
			t.Fatal(err)
		}
	}
	return fp.Name()
}

func assertFileSame(t *testing.T, filePathExpected string, filePathActual string) {
	expected, err := os.ReadFile(filePathExpected)
	assert.NoError(t, err)

	actual, err := os.ReadFile(filePathActual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func Test_updateRelPath(t *testing.T) {
	type args struct {
		line    string
		srcDir  string
		destDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				line:    `  source = "./../modules"`,
				srcDir:  `a/b/c`,
				destDir: `a/x/y`,
			},
			want: `  source = "../../b/modules"`,
		},
		{
			args: args{
				line:    `  source = "./modules"`,
				srcDir:  `a/b/c`,
				destDir: `a/x/y`,
			},
			want: `  source = "../../b/c/modules"`,
		},
		{
			args: args{
				line:    `  source = "git.com/modules"`,
				srcDir:  `a/b/c`,
				destDir: `a/x/y`,
			},
			want: `  source = "git.com/modules"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateRelPath(tt.args.line, tt.args.srcDir, tt.args.destDir)
			assert.Equal(t, tt.want, got)
		})
	}
}
