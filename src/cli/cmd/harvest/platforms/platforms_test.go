// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platforms

import (
	"context"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/git"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	terrpb "github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_parseYAML(t *testing.T) {
	validYAML := []byte(`
components:
  - id: ComponentA
    title: TypeA
  - id: ComponentB
    title: TypeB
`)

	expectedData := YAMLData{
		Components: []Component{
			{ID: "ComponentA", Title: "TypeA"},
			{ID: "ComponentB", Title: "TypeB"},
		},
	}

	data, err := parseYAML(validYAML)
	require.NoError(t, err, "Expected no error when parsing valid YAML")
	assert.Equal(t, expectedData, data, "Parsed data should match expected data")

	invalidYAML := []byte("invalid: yaml: content")

	data, err = parseYAML(invalidYAML)
	require.Error(t, err, "Expected an error when parsing invalid YAML")
	assert.Empty(t, data, "Parsed data should be empty when YAML is invalid")
}

func Test_createDBPlatform(t *testing.T) {
	platform := Platform{
		Title:       "Sample Title",
		Description: "Sample Description",
		RepoURL:     "https://github.com/sample/repo",
		RepoDir:     "sample-directory",
	}
	commitSHA := "sampleSHA"
	revision := Revision{
		Label: "ref",
		Type:  "branch",
	}

	dbPlatform := createDBPlatform(platform, commitSHA, revision)

	assert.Equal(t, platform.Title, dbPlatform.Title, "Title should match")
	assert.Equal(t, platform.Description, dbPlatform.Description, "Description should match")
	assert.Equal(t, platform.RepoURL, dbPlatform.RepoURL, "RepoURL should match")
	assert.Equal(t, platform.RepoDir, dbPlatform.RepoDirectory, "RepoDirectory should match")
	assert.Equal(t, commitSHA, dbPlatform.CommitSHA, "CommitSHA should match")
	assert.Equal(t, revision.Label, dbPlatform.RefLabel, "RefLabel should match")
	assert.Equal(t, terrpb.GitLabelEnum(1), dbPlatform.LabelType, "LabelType should match (assuming 'branch' corresponds to 1)")
}

func Test_getOwnerRepoRef(t *testing.T) {
	platform := Platform{
		RepoURL: "https://github.com/owner/repo",
		Revisions: []Revision{
			{Label: "main"},
			{Label: "feature/branch"},
		},
	}

	owner, repo, ref, err := getOwnerRepoRef(platform)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "owner", owner, "Owner should match")
	assert.Equal(t, "repo", repo, "Repo should match")
	assert.Equal(t, "main", ref, "Ref should match")
}

func TestGitLabelEnumFromYAML(t *testing.T) {
	testCases := []struct {
		language       string
		expectedResult int32
	}{
		{"branch", int32(terrariumpb.GitLabelEnum_label_branch)},
		{"tag", int32(terrariumpb.GitLabelEnum_label_tag)},
		{"commit", int32(terrariumpb.GitLabelEnum_label_commit)},
		{"unknown", int32(terrariumpb.GitLabelEnum_label_no)},
	}

	for _, tc := range testCases {
		t.Run(tc.language, func(t *testing.T) {
			result := GitLabelEnumFromYAML(tc.language)
			if result != tc.expectedResult {
				t.Errorf("Expected %d for language '%s', but got %d", tc.expectedResult, tc.language, result)
			}
		})
	}
}

func Test_isPlatformYAML(t *testing.T) {
	// Test when the filename is "platform.yaml"
	t.Run("Test with platform.yaml", func(t *testing.T) {
		result := isPlatformYAML("platform.yaml")
		if !result {
			t.Errorf(`isPlatformYAML("platform.yaml") = %v; want true`, result)
		}
	})

	// Test when the filename is "platform.yml"
	t.Run("Test with platform.yml", func(t *testing.T) {
		result := isPlatformYAML("platform.yml")
		if !result {
			t.Errorf(`isPlatformYAML("platform.yml") = %v; want true`, result)
		}
	})

	// Test with other filenames
	t.Run("Test with other filename", func(t *testing.T) {
		result := isPlatformYAML("otherfile.txt")
		if result {
			t.Errorf(`isPlatformYAML("otherfile.txt") = %v; want false`, result)
		}
	})
}

func Test_findPlatformYAML(t *testing.T) {
	tests := []struct {
		name          string
		directoryPath string
		extensions    map[string]bool
		wantPath      string
		wantErr       error
	}{
		{
			name:          "Valid directory with platform.yaml",
			directoryPath: "../../../../../examples/farm/platform/",
			extensions: map[string]bool{
				".yaml": true,
				".yml":  true,
			},
			wantPath: "../../../../../examples/farm/platform/platform.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundPath, err := findPlatformYAML(tt.directoryPath)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPath, foundPath)
			}
		})
	}
}

func TestCmd(t *testing.T) {
	testSetup := clitesting.CLITest{
		CmdToTest: NewCmd,
		SetupTest: func(ctx context.Context, t *testing.T) {
			t.Setenv("TR_LOG_LEVEL", "error")
			config.LoadDefaults()
		},
	}

	testSetup.RunTests(t, []clitesting.CLITestCase{
		{
			Name: "error connecting to db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				config.SetDBMocks(nil)
			},
			Args:     []string{},
			WantErr:  true,
			ExpError: "error connecting to the database: mocked err: connection failed",
		},
		{
			Name: "error connecting to db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				config.SetDBMocks(&mocks.DB{})
			},
			Args:     []string{"-d", "./non-existing"},
			WantErr:  true,
			ExpError: "error parsing platform YAML: stat ./non-existing: no such file or directory",
		},
	})
}

func Test_Component_Init(t *testing.T) {
	component := &Component{
		Inputs:  &jsonschema.Node{Type: "previous_type"},
		Outputs: &jsonschema.Node{Type: "previous_type"},
	}

	component.Init()

	assert.Equal(t, "object", component.Inputs.Type)
	assert.Equal(t, "object", component.Outputs.Type)
}

func TestTerrarium_Init(t *testing.T) {
	terrarium := Terrarium{
		Components{
			{
				Title:   "Test Component 1",
				Inputs:  &jsonschema.Node{Type: "previous_type"},
				Outputs: &jsonschema.Node{Type: "previous_type"},
			},
			{
				Title:   "Test Component 2",
				Inputs:  &jsonschema.Node{Type: "previous_type"},
				Outputs: &jsonschema.Node{Type: "previous_type"},
			},
		},
	}

	terrarium.Init()

	for _, component := range terrarium.Terrarium {
		assert.Equal(t, "object", component.Inputs.Type)
		assert.Equal(t, "object", component.Outputs.Type)
	}
}

func Test_harvestPlatforms(t *testing.T) {
	q := []db.DependencyResult{
		{
			DependencyID: uuid.New(),
			InterfaceID:  "test_id",
			Name:         "test_name",
			Schema:       nil,
			Computed:     false,
		},
	}
	tests := []struct {
		name    string
		dirPath string
		mockDB  func(*mocks.DB)
		mockGit func(*git.Git)
		wantErr bool
	}{
		{
			name:    "success with valid YAML file",
			dirPath: "./testdata/",
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("CreatePlatform", mock.Anything).Return(uuid.New(), nil).Once()
				dbMocks.On("Fetchdeps", mock.Anything).Return(q, nil).Once()
			},
		},
	}
	for _, tt := range tests {
		config.LoadDefaults()
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			if tt.mockDB != nil {
				tt.mockDB(dbMocks)
			}

			err := harvestPlatforms(dbMocks, tt.dirPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			dbMocks.AssertExpectations(t)
		})
	}
}

func Test_readPlatformYAML(t *testing.T) {
	tests := []struct {
		name    string
		dirPath string
		wantErr bool
	}{
		{
			name:    "success with dir",
			dirPath: "./",
		},
		{
			name:    "success with file",
			dirPath: "./testdata/platform.yaml",
		},
	}
	for _, tt := range tests {
		config.LoadDefaults()
		t.Run(tt.name, func(t *testing.T) {
			_, err := readPlatformYAML(tt.dirPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_findTerrariumYaml(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		gl        []*github.RepositoryContent
		owner     string
		repo      string
		reference string
		dirPath   string
		mockDB    func(*mocks.DB)
		wantErr   bool
	}{
		{
			name:      "success",
			owner:     "owner",
			repo:      "repo",
			reference: "ref",
			dirPath:   "./",
			gl: []*github.RepositoryContent{
				{
					Name: github.String("platform.yaml"),
					Type: github.String("dir"),
					Path: github.String("./testdata"),
				},
			},
		},
	}
	for _, tt := range tests {
		config.LoadDefaults()
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			if tt.mockDB != nil {
				tt.mockDB(dbMocks)
			}

			_, err := findTerrariumYaml(ctx, tt.gl, tt.owner, tt.repo, tt.reference, tt.dirPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			dbMocks.AssertExpectations(t)
		})
	}
}

func Test_compareYAMLWithSQLData(t *testing.T) {
	tests := []struct {
		name        string
		queryResult []db.DependencyResult
		c           Component
		u           uuid.UUID
		mockDB      func(*mocks.DB)
		wantErr     bool
	}{
		{
			name: "success",
			queryResult: []db.DependencyResult{
				{
					DependencyID: uuid.New(),
					InterfaceID:  "testID",
					Name:         "test_name",
					Schema:       nil,
					Computed:     false,
				},
			},
			c: Component{
				ID:          "testID",
				Title:       "testTitle",
				Description: "testDesc",
				Inputs:      nil,
				Outputs:     nil,
			},
			u: uuid.New(),
			mockDB: func(dbMocks *mocks.DB) {
				dbMocks.On("CreatePlatformComponents", mock.Anything).Return(uuid.New(), nil).Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		config.LoadDefaults()
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			if tt.mockDB != nil {
				tt.mockDB(dbMocks)
			}

			err := compareYAMLWithSQLData(dbMocks, tt.c, tt.queryResult, tt.u)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			dbMocks.AssertExpectations(t)
		})
	}
}
