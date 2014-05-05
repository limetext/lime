// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// This file was generated as part of a build step and shouldn't be manually modified

package sublime

import (
	"fmt"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	"github.com/quarnster/util/text"
)

var (
	_ = backend.View{}
	_ = text.Region{}
	_ = fmt.Errorf
)

var _editClass = py.Class{
	Name:    "sublime.Edit",
	Pointer: (*Edit)(nil),
}

type Edit struct {
	py.BaseObject
	data *backend.Edit
}

func (o *Edit) PyInit(args *py.Tuple, kwds *py.Dict) error {
	return fmt.Errorf("Can't initialize type Edit")
}
func (o *Edit) PyStr() string {
	return o.data.String()
}
