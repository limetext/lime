// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	// The ToggleSettingCommand toggles the value of a setting,
	// making it false when it was true or true when it was false.
	ToggleSettingCommand struct {
		BypassUndoCommand
		Setting string
	}

	// The SetSettingCommand set the value of a setting.
	SetSettingCommand struct {
		BypassUndoCommand
		Setting string
		Value   interface{}
	}
)

func (c *ToggleSettingCommand) Run(v *View, e *Edit) error {
	setting := c.Setting
	prev, boolean := v.Settings().Get(setting, false).(bool)
	// if the setting was non-boolean, it is set to true, else it is toggled
	v.Settings().Set(setting, !boolean || !prev)
	return nil
}

func (c *SetSettingCommand) Run(v *View, e *Edit) error {
	setting := c.Setting
	v.Settings().Set(setting, c.Value)
	return nil
}

func init() {
	register([]Command{
		&ToggleSettingCommand{},
		&SetSettingCommand{},
	})
}
