package terrariumsrv

import (
	"testing"

	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultPage(t *testing.T) {
	testCases := []struct {
		name         string
		page         *terrariumpb.Page
		expectedSize int32
	}{
		{
			name:         "Nil Page",
			page:         nil,
			expectedSize: 100,
		},
		{
			name:         "Non-Nil Page",
			page:         &terrariumpb.Page{Size: 50},
			expectedSize: 50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newPage := setDefaultPage(tc.page)

			assert.NotNil(t, newPage)
			assert.Equal(t, tc.expectedSize, newPage.Size)
		})
	}
}
