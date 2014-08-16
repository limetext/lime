// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"
	"testing"
)

type (
	DummyApplicationCommand struct {
		DefaultCommand
	}

	DummyWindowCommand struct {
		DefaultCommand
	}

	DummyTextCommand struct {
		DefaultCommand
	}
)

func (c *DummyApplicationCommand) Run() error {
	return fmt.Errorf("Ran")
}

func (c *DummyApplicationCommand) IsChecked() bool {
	return false
}

func (c *DummyWindowCommand) Run(w *Window) error {
	return fmt.Errorf("Ran")
}

func (c *DummyTextCommand) Run(v *View, e *Edit) error {
	return fmt.Errorf("Ran")
}

func TestPascalCaseToSnakeCase(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			"TestString",
			"test_string",
		},
		{
			"Teststring",
			"teststring",
		},
	}

	for i, test := range tests {
		out := PascalCaseToSnakeCase(test.in)

		if out != test.out {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.out, out)
		}
	}
}

func TestRegisterApplicationCommand(t *testing.T) {
	name := "app_test_command"
	ac := DummyApplicationCommand{}
	ch := GetEditor().CommandHandler()

	err := ch.Register(name, &ac)

	if err != nil {
		t.Errorf("Got error while registering: %s", err)
	}

	err = ch.RunApplicationCommand(name, Args{})

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}

func TestRegisterWindowCommand(t *testing.T) {
	var w Window

	name := "wnd_test_command"
	wc := DummyWindowCommand{}
	ch := GetEditor().CommandHandler()

	err := ch.Register(name, &wc)

	if err != nil {
		t.Errorf("Got error while registering: %s", err)
	}

	err = ch.RunWindowCommand(&w, name, Args{})

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}

func TestRegisterTextCommand(t *testing.T) {
	ed := GetEditor()

	name := "text_test_command"
	tc := DummyTextCommand{}
	ch := ed.CommandHandler()

	err := ch.Register(name, &tc)

	if err != nil {
		t.Errorf("Got error while registering: %s", err)
	}

	v := ed.NewWindow().NewFile()
	err = ch.RunTextCommand(v, name, Args{})

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}
