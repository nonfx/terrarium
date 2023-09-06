// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import "github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"

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
