// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"testing"
)

func TestPlatform_GetCondition(t *testing.T) {
	pl := &Platform{
		Title: "TestPlatform",
	}

	condition := pl.GetCondition()

	// Check if the condition has the same Name as the original Platform
	if condition.(*Platform).Title != pl.Title {
		t.Errorf("GetCondition() returned an incorrect condition. Expected Name: %s, Got Name: %s", pl.Title, condition.(*Platform).Title)
	}
}
