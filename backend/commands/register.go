// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"code.google.com/p/log4go"
	"lime/backend"
)

type cmd struct {
	name string
	cmd  backend.Command
}

func register(cmds []cmd) {
	e := backend.GetEditor()
	for i := range cmds {
		if err := e.CommandHandler().Register(cmds[i].name, cmds[i].cmd); err != nil {
			log4go.Error("Failed to register command %s: %s", cmds[i].name, err)
		}
	}
}
