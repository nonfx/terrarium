// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package config

import "github.com/cldcvr/terrarium/src/pkg/confighelper"

func GitToken() string {
	return confighelper.MustGetString("github.token")
}
