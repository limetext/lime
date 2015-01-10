// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"testing"
)

func TestRunApplicationCommand(t *testing.T) {
	nopApplicationCommand := NopApplicationCommand{}

	if nopApplicationCommand.Run() != nil {
		t.Error("No op application command running returns not nil")
	}

}

func TestRunNopWindowCommand(t *testing.T) {
	nopWindowCommand := NopWindowCommand{}

	if nopWindowCommand.Run(&Window{}) != nil {
		t.Error("No op window command running returns not nil")
	}
}

func TestRunNopTextCommand(t *testing.T) {
	nopTextCommand := NopTextCommand{}

	if nopTextCommand.Run(&View{}, &Edit{}) != nil {
		t.Error("No op text command running returns not nil")
	}
}
