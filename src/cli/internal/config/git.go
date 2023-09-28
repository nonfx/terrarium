// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package config

import "github.com/cldcvr/terrarium/src/pkg/confighelper"

func GitUsername() string {
	return confighelper.MustGetString("github.username")
}

func GitPassword() string {
	return confighelper.MustGetString("github.password")
}
