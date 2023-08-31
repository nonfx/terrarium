// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package commander

import "os/exec"

//go:generate mockery --all

type Commander interface {
	Run(*exec.Cmd) error
}
