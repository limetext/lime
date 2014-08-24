// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"fmt"
	"testing"
)

type DummyApplicationCommand struct {
	DefaultCommand
}

func (c *DummyApplicationCommand) Run() error {
	return fmt.Errorf("Ran")
}

func (c *DummyApplicationCommand) IsChecked() bool {
	return false
}

func TestRegister(t *testing.T) {
	ed := GetEditor()

	name := "dummy"

	register([]cmd {
		{name, &DummyApplicationCommand{}},
	})

	err := ed.CommandHandler().RunApplicationCommand(name, nil)

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}
