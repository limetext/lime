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

func TestDefaultName(t *testing.T) {
	n := DefaultName(&DummyApplicationCommand{})

	if n != "dummy_application" {
		t.Errorf("Expected %s, but got %s", "dummy_application", n)
	}
}

func TestRegisterAndRunApplicationCommand(t *testing.T) {
	name := "app_test"
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

func TestRegisterAndRunWindowCommand(t *testing.T) {
	ed := GetEditor()

	name := "wnd_test"
	wc := DummyWindowCommand{}
	ch := ed.CommandHandler()

	err := ch.Register(name, &wc)

	if err != nil {
		t.Errorf("Got error while registering: %s", err)
	}

	w := ed.NewWindow()
	defer w.Close()

	err = ch.RunWindowCommand(w, name, Args{})

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}

func TestRegisterAndRunTextCommand(t *testing.T) {
	ed := GetEditor()

	name := "text_test"
	tc := DummyTextCommand{}
	ch := ed.CommandHandler()

	err := ch.Register(name, &tc)

	if err != nil {
		t.Errorf("Got error while registering: %s", err)
	}

	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	err = ch.RunTextCommand(v, name, Args{})

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}

func TestRegisterAndRunDefaultNamedCommand(t *testing.T) {
	ac := DummyApplicationCommand{}
	ch := GetEditor().CommandHandler()

	err := ch.RegisterWithDefault(&ac)

	if err != nil {
		t.Errorf("Got error while registering: %s", err)
	}

	name := DefaultName(&ac)
	err = ch.RunApplicationCommand(name, Args{})

	if err == nil {
		t.Errorf("Expected %s to run, but it didn't", name)
	} else if err.Error() != "Ran" {
		t.Errorf("Expected %s to run, but it got an error: %v", name, err)
	}
}

func TestUnregisterAndRunApplicationCommand(t *testing.T) {
	name := "app_test_unregister"
	ac := DummyApplicationCommand{}
	ch := GetEditor().CommandHandler()

	_ = ch.Register(name, &ac)
	err := ch.Unregister(name)

	if err != nil {
		t.Errorf("Got error while unregistering: %s", err)
	}

	err = ch.RunApplicationCommand(name, Args{})

	if err != nil {
		t.Errorf("Expected %s not to run, but it did", name)
		if err.Error() != "Ran" {
			t.Errorf("Expected %s not to run, but it got an error: %v", name, err)
		}
	}
}

func TestUnregisterAndRunWindowCommand(t *testing.T) {
	ed := GetEditor()

	name := "wnd_test_unregister"
	wc := DummyWindowCommand{}
	ch := ed.CommandHandler()

	_ = ch.Register(name, &wc)
	err := ch.Unregister(name)

	if err != nil {
		t.Errorf("Got error while unregistering: %s", err)
	}

	w := ed.NewWindow()
	defer w.Close()

	err = ch.RunWindowCommand(w, name, Args{})

	if err != nil {
		t.Errorf("Expected %s not to run, but it did", name)
		if err.Error() != "Ran" {
			t.Errorf("Expected %s not to run, but it got an error: %v", name, err)
		}
	}
}

func TestUnregisterAndRunTextCommand(t *testing.T) {
	ed := GetEditor()

	name := "text_test"
	tc := DummyTextCommand{}
	ch := ed.CommandHandler()

	_ = ch.Register(name, &tc)
	err := ch.Unregister(name)

	if err != nil {
		t.Errorf("Got error while unregistering: %s", err)
	}

	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	err = ch.RunTextCommand(v, name, Args{})

	if err != nil {
		t.Errorf("Expected %s not to run, but it did", name)
		if err.Error() != "Ran" {
			t.Errorf("Expected %s not to run, but it got an error: %v", name, err)
		}
	}
}
