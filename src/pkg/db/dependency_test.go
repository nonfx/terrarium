// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockgDB struct {
	mock.Mock
}

type MockDbChain struct {
	mock.Mock
}

func (m *MockgDB) g() *MockDbChain {
	return &MockDbChain{}
}

func (chain *MockDbChain) Where(query interface{}, args ...interface{}) *MockDbChain {
	return chain
}

func (chain *MockDbChain) Preload(column string) *MockDbChain {
	return chain
}

func (chain *MockDbChain) First(dest interface{}) *MockDbChain {
	args := chain.Called(dest)
	return args.Get(0).(*MockDbChain)
}

func (chain *MockDbChain) Error() error {
	args := chain.Called()
	return args.Error(0)
}

func (m *MockgDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	args := m.Called(e)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockgDB) FetchDependencyByInterfaceID(interfaceID string) (*terrariumpb.Dependency, error) {
	args := m.Called(interfaceID)
	return args.Get(0).(*terrariumpb.Dependency), args.Error(1)
}

func TestCreateDependencyInterface(t *testing.T) {
	mockDB := new(MockgDB)
	dep := &Dependency{InterfaceID: "test_id"}

	// Mock db call.
	uuidVal, _ := uuid.NewUUID()
	mockDB.On("CreateDependencyInterface", dep).Return(uuidVal, nil)

	_, err := mockDB.CreateDependencyInterface(dep)
	assert.Nil(t, err)
}

func TestToProto(t *testing.T) {
	dep := &Dependency{
		InterfaceID: "test_id",
		Inputs:      &jsonschema.Node{Type: "string"},
		Outputs:     &jsonschema.Node{Type: "string"},
	}

	protoDep, err := dep.ToProto()
	assert.Nil(t, err)
	assert.Equal(t, "test_id", protoDep.InterfaceId)
}

func TestFetchDependencyByInterfaceID(t *testing.T) {
	mockDB := new(MockgDB)
	interfaceID := "test_id"

	// Mock db call.
	dep := &terrariumpb.Dependency{InterfaceId: "test_id"}
	mockDB.On("FetchDependencyByInterfaceID", interfaceID).Return(dep, nil)

	protoDep, err := mockDB.FetchDependencyByInterfaceID(interfaceID)
	assert.Nil(t, err)
	assert.Equal(t, "test_id", protoDep.InterfaceId)
}
