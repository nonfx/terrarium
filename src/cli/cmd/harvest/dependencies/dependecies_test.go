package dependencies

import (
	"io/fs"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockFS is a mock implementation of the fs.FS interface.
type MockFS struct {
	mock.Mock
	files map[string]string
}

type Taxonomy struct {
	ID uuid.UUID
}

type MockDB struct {
	mock.Mock
}

func NewMockFS(files map[string]string) fs.FS {
	return &MockFS{files: files}
}

func (m *MockFS) Open(name string) (fs.File, error) {
	if content, ok := m.files[name]; ok {
		return NewMockFile(content), nil
	}
	return nil, fs.ErrNotExist
}

type MockFile struct {
	mock.Mock
	content string
}

func NewMockFile(content string) fs.File {
	return &MockFile{content: content}
}

func (m *MockFile) Stat() (fs.FileInfo, error) {
	return &MockFileInfo{}, nil
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	copy(p, []byte(m.content))
	return len(m.content), nil
}

func (m *MockFile) Close() error {
	return nil
}

type MockFileInfo struct {
	mock.Mock
	content string
}

func (m *MockFileInfo) Name() string       { return "mockfile" }
func (m *MockFileInfo) Size() int64        { return int64(len(m.content)) }
func (m *MockFileInfo) Mode() fs.FileMode  { return fs.ModePerm }
func (m *MockFileInfo) ModTime() time.Time { return time.Now() } // Implement ModTime
func (m *MockFileInfo) IsDir() bool        { return false }
func (m *MockFileInfo) Sys() interface{}   { return nil }

// Example usage in a test:
// func TestProcessYAMLFiles(t *testing.T) {
// 	mockDB := new(MockDB)
// 	mockDB.On("GetTaxonomyByFieldName", mock.Anything, mock.Anything).Return(Taxonomy{}, nil)

// 	// data := []byte(`...`) // Replace with your test YAML data

// 	// mockFS := NewMockFS(map[string]string{
// 	// 	"file1.yaml": string(data),
// 	// 	// Add more files as needed
// 	// })

// 	err := processYAMLFiles("testdir")
// 	assert.NoError(t, err)
// }
