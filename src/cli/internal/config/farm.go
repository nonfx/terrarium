// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package config

import "github.com/cldcvr/terrarium/src/pkg/confighelper"

// FarmDefault link of the farm repo to use by default
func FarmDefault() string {
	return confighelper.MustGetString("farm.default")
}

// FarmVersion version of the farm repo to use by default
func FarmVersion() string {
	return confighelper.MustGetString("farm.version")
}
