// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package config

import "github.com/cldcvr/terrarium/src/pkg/confighelper"

// FarmDefault link to the farm repo being utilized
func FarmDefault() string {
	return confighelper.MustGetString("farm.default")
}
