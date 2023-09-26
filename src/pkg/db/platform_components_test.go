// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"testing"

	"github.com/google/uuid"
)

func TestPlatformComponents_GetCondition(t *testing.T) {
	// Create a PlatformComponents instance with specific PlatformID and DependencyID
	pc := &PlatformComponents{
		PlatformID:   uuid.MustParse("f1b5ab7e-9d10-4f21-9752-0d3984aa2b0e"),
		DependencyID: uuid.MustParse("40b0be65-3ac5-4a31-8bc6-125734c3db4e"),
	}

	// Call the GetCondition method
	condition := pc.GetCondition()

	// Check if the condition has the same PlatformID and DependencyID as the original PlatformComponents
	conditionPC, ok := condition.(*PlatformComponents)
	if !ok {
		t.Errorf("GetCondition() did not return a PlatformComponents condition")
		return
	}

	if conditionPC.PlatformID != pc.PlatformID || conditionPC.DependencyID != pc.DependencyID {
		t.Errorf("GetCondition() returned an incorrect condition. Expected PlatformID: %s, Got PlatformID: %s, Expected DependencyID: %s, Got DependencyID: %s", pc.PlatformID, conditionPC.PlatformID, pc.DependencyID, conditionPC.DependencyID)
	}
}
