// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"testing"
)

func TestDefaultCommand(t *testing.T) {
	dc := DefaultCommand{}

	if dc.BypassUndo() != false {
		t.Errorf("Expected BypassUndo to return false, but got %v", dc.BypassUndo())
	}

	if dc.IsEnabled() != true {
		t.Errorf("Expected IsEnabled to return true, but got %v", dc.IsEnabled())
	}

	if dc.IsVisible() != true {
		t.Errorf("Expected IsVisible to return true, but got %v", dc.IsVisible())
	}

	if dc.Description() != "TODO" {
		t.Errorf("Expected Description to return \"TODO\", but got %v", dc.Description())
	}
}

func TestBypassUndoCommand(t *testing.T) {
	bc := BypassUndoCommand{}

	if bc.BypassUndo() != true {
		t.Errorf("Expected BypassUndo to return true, but got %v", bc.BypassUndo())
	}

	if bc.IsEnabled() != true {
		t.Errorf("Expected IsEnabled to return true, but got %v", bc.IsEnabled())
	}

	if bc.IsVisible() != true {
		t.Errorf("Expected IsVisible to return true, but got %v", bc.IsVisible())
	}

	if bc.Description() != "TODO" {
		t.Errorf("Expected Description to return \"TODO\", but got %v", bc.Description())
	}
}
