//go:build dbtest
// +build dbtest

package db

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_gDB_FindOutputMappingsByModuleID(t *testing.T) {
	db, err := Connect()
	require.NoError(t, err)

	db = (*gDB)(db.(*gDB).g().Debug())

	tests := []struct {
		name       string
		id         uuid.UUID
		wantResult *TFModule
		wantErr    bool
	}{
		{
			id: uuid.MustParse("33be72db-fa09-4c7b-a66d-24b8a2e1ae93"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := db.FindOutputMappingsByModuleID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("gDB.FindOutputMappingsByModuleID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("gDB.FindOutputMappingsByModuleID() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
