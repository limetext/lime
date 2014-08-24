// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	. "github.com/limetext/lime/backend"
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

func TestRegisterByName(t *testing.T) {
	ed := GetEditor()

	name := "dummy"

	registerByName([]namedCmd{
		{name, &DummyApplicationCommand{}},
	})

	err := ed.CommandHandler().RunApplicationCommand(name, nil)

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}

func TestRegister(t *testing.T) {
	ed := GetEditor()
	ac := &DummyApplicationCommand{}

	register([]Command{
		ac,
	})

	name := DefaultName(ac)
	err := ed.CommandHandler().RunApplicationCommand(name, nil)

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}
