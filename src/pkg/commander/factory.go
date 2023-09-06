// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build !mock
// +build !mock

package commander

func GetCommander() Commander {
	return &osExec{}
}
