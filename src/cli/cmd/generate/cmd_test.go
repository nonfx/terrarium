package generate

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCmd(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectOut   string
		expectError bool
	}{
		{
			name:        "No components provided",
			args:        []string{},
			expectError: true,
		},
		{
			args:      []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium"},
			expectOut: "Successfully pulled 9 of 15 terraform blocks at: ./testdata/.terrarium\n",
		},
	}

	os.RemoveAll("./testdata/.terrarium")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetCmd()
			out := &strings.Builder{}
			cmd.SetOut(out)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectOut, out.String())
			}
		})
	}
}
