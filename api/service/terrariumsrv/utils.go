package terrariumsrv

import "github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"

// setDefaultPage returns new page object with default values set. doesn't mutate the given page object.
func setDefaultPage(page *terrariumpb.Page) (newPage *terrariumpb.Page) {
	if page != nil {
		newPage = page
	} else {
		newPage = &terrariumpb.Page{}
	}

	if newPage.Size == 0 {
		newPage.Size = 100
	}

	return
}
