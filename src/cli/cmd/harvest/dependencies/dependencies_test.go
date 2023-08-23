package dependencies

import (
	"fmt"
	"os"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_processYAMLData(t *testing.T) {
	tests := []struct {
		name     string
		yamlData []byte
		mockDB   func(*mocks.DB)
		wantErr  bool
	}{
		{
			name: "success with taxonomy",
			yamlData: []byte(`
  dependency-interfaces:
  - id: interface1
    taxonomy: storage/database/rdbms
    title: RDBMS
    description: Relational Database Management System
    inputs:
      type: object
      properties:
        port:
          title: Port
          description: The port number on which the server should listen.
          type: number
          default: 80
    outputs:
      properties:
        host:
          title: Host
          description: The host address of the web server.
          type: string
    `),
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("GetTaxonomyByFieldName", mock.Anything, mock.Anything).Return(
					db.Taxonomy{Level1: "level1", Level2: "level2", Level3: "level3"}, nil).Once()
				dbMocks.On("GetTaxonomyByFieldName", mock.Anything, mock.Anything).Return(
					db.Taxonomy{Level1: "level1", Level2: "level2", Level3: "level3"}, nil).Once()
				dbMocks.On("GetTaxonomyByFieldName", mock.Anything, mock.Anything).Return(
					db.Taxonomy{Level1: "level1", Level2: "level2", Level3: "level3"}, nil).Once()
				dbMocks.On("CreateDependencyInterface", mock.Anything).Return(uuid.New(), nil).Once()
			},
		},
		{
			name: "failure due to unmarshal error",
			yamlData: []byte(`
				dependency-interfaces:
				  - invalid
				`),
			wantErr: true,
		},
		{
			name: "success without taxonomy",
			yamlData: []byte(`
dependency-interfaces:
  - id: dep2
    title: Dependency 2
    description: Description for Dependency 2
    inputs: {}
    outputs: {}
`),
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("CreateDependencyInterface", mock.Anything).Return(uuid.New(), nil).Once()
			},
		},
		{
			name: "failure with taxonomy",
			yamlData: []byte(`
dependency-interfaces:
  - id: dep3
    taxonomy: storage/database
    title: Dependency 3
    description: Description for Dependency 3
    inputs: {}
    outputs: {}
`),
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("GetTaxonomyByFieldName", mock.Anything, "storage").Return(
					db.Taxonomy{}, fmt.Errorf("mocked error")).Once()
			},
			wantErr: true,
		},
		{
			name: "failure on CreateDependencyInterface",
			yamlData: []byte(`
dependency-interfaces:
  - id: dep4
    title: Dependency 4
    description: Description for Dependency 4
    inputs: {}
    outputs: {}
`),
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("CreateDependencyInterface", mock.Anything).Return(uuid.Nil, fmt.Errorf("mocked error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			if tt.mockDB != nil {
				tt.mockDB(dbMocks)
			}

			err := processYAMLData(dbMocks, "dummy.yaml", []byte(tt.yamlData))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			dbMocks.AssertExpectations(t)
		})
	}
}

func Test_cmdRunE(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	config.LoadDefaults()
	os.Setenv("TR_DB_TYPE", "mock")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "DBConnect Error",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			wantErr: true,
		},
		{
			name: "ProcessYAMLFiles Error",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmdRunE(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("cmdRunE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockOS struct {
	mock.Mock
}

func (m *mockOS) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

type OSReader interface {
	ReadFile(filename string) ([]byte, error)
}

func Test_processYAMLFiles(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		mockDB    func(*mocks.DB)
		mockOS    func() *mockOS
		wantErr   bool
	}{
		{
			name:      "success with valid YAML file",
			directory: "testdata/success",
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("GetTaxonomyByFieldName", "level1", mock.Anything).Return(
					db.Taxonomy{Level1: "level1", Level2: "level2", Level3: "level3"}, nil).Once()
				dbMocks.On("GetTaxonomyByFieldName", "level2", mock.Anything).Return(
					db.Taxonomy{Level1: "level1", Level2: "level2", Level3: "level3"}, nil).Once()
				dbMocks.On("GetTaxonomyByFieldName", "level3", mock.Anything).Return(
					db.Taxonomy{Level1: "level1", Level2: "level2", Level3: "level3"}, nil).Once()
				dbMocks.On("CreateDependencyInterface", mock.Anything).Return(uuid.New(), nil).Once()
			},
		},
		{
			name:      "failure processing YAML data",
			directory: "testdata/failure",
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("CreateDependencyInterface", mock.Anything).Return(uuid.New(), fmt.Errorf("mocked error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			osMock := &mockOS{}
			if tt.mockDB != nil {
				tt.mockDB(dbMocks)
			}

			err := processYAMLFiles(dbMocks, tt.directory)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			dbMocks.AssertExpectations(t)
			osMock.AssertExpectations(t)
		})
	}
}

func TestGetCmd(t *testing.T) {
	cmd := GetCmd()
	assert.NotNil(t, cmd)
}
