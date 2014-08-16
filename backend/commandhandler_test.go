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
)

func (c *DummyApplicationCommand) Run() error {
	return fmt.Errorf("Ran")
}

func (c *DummyApplicationCommand) IsChecked() bool {
	return false
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
