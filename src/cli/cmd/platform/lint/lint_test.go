package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lintPlatform(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid platform",
			args: args{
				dir: "testdata/valid-terraform-1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lintPlatform(tt.args.dir)
			assert.NoError(t, err)
		})
	}
}
