// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	. "github.com/limetext/lime/backend"
)

type (
	OpenFileCommand struct {
		DefaultCommand
		Path string
	}
)

func (o *OpenFileCommand) Set(v interface{}) error {
	str, ok := v.(string)
	if ok {
		o.Path = str
	} else {
		return fmt.Errorf("%s is of unknown type", v)
	}
	return nil
}

func (o *OpenFileCommand) Run(w *Window) error {
	ed := GetEditor()
	ed.ActiveWindow().OpenFile(o.Path, 0)
	return nil
}

func init() {
	register([]Command{
		&OpenFileCommand{},
	})
}
