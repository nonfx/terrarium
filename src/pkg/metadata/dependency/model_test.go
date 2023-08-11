package dependency

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		expected *File
		err      bool
	}{
		{
			name: "Valid YAML data",
			data: []byte(`dependency-interfaces:
  - id: postgres
    title: PostgreSQL Database
    description: A relational database management system using SQL.
    inputs:
      properties:
          db_name:
              title: Database name
              description: The name provided here may get prefix and suffix based
              type: string
              default: random
          version:
              title: Engine version
              description: Version of the PostgreSQL engine to use
              type: string
              default: "11"
    outputs:
      properties:
          host:
              title: Host
              description: The host address of the PostgreSQL server.
              type: string
          password:
              title: Password
              description: The password for accessing the PostgreSQL database.
              type: string
          port:
              title: Port
              description: The port number on which the PostgreSQL server is listening.
              type: number
          username:
              title: Username
              description: The username for accessing the PostgreSQL database.
              type: string
`),
			expected: &File{
				DependencyInterfaces: Interfaces{
					Interface{
						ID:          "postgres",
						Title:       "PostgreSQL Database",
						Description: "A relational database management system using SQL.",
						Inputs: &jsonschema.Node{
							Type: "object",
							Properties: map[string]*jsonschema.Node{
								"db_name": {
									Title:       "Database name",
									Description: "The name provided here may get prefix and suffix based",
									Type:        "string",
									Default:     "random",
								},
								"version": {
									Title:       "Engine version",
									Description: "Version of the PostgreSQL engine to use",
									Type:        "string",
									Default:     "11",
								},
							},
						},

						Outputs: &jsonschema.Node{
							Type: "object",
							Properties: map[string]*jsonschema.Node{
								"host": {
									Title:       "Host",
									Description: "The host address of the PostgreSQL server.",
									Type:        "string",
								},
								"password": {
									Title:       "Password",
									Description: "The password for accessing the PostgreSQL database.",
									Type:        "string",
								},
								"port": {
									Title:       "Port",
									Description: "The port number on which the PostgreSQL server is listening.",
									Type:        "number",
								},
								"username": {
									Title:       "Username",
									Description: "The username for accessing the PostgreSQL database.",
									Type:        "string",
								},
							},
						},
					},
				},
			},
			err: false,
		},
		{
			name:     "Invalid YAML data",
			data:     []byte(`invalid: yaml: data`),
			expected: nil,
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewFile(tc.data)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
