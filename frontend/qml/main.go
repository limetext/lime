// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	_ "github.com/limetext/lime/backend/commands"
	"github.com/limetext/lime/backend/keys"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/render"
	"github.com/limetext/lime/backend/textmate"
	"gopkg.in/qml.v1"
	"runtime"
	"sync"
)

const (
	qmlMainFile = "main.qml"
	qmlViewFile = "LimeView.qml"
)

var (
	limeViewComponent qml.Object
	scheme            *textmate.Theme
	blink             bool
)

const (
	console_height  = 20
	render_chan_len = 2

	// http://qt-project.org/doc/qt-5.1/qtcore/qt.html#KeyboardModifier-enum
	shift_mod  = 0x02000000
	ctrl_mod   = 0x04000000
	alt_mod    = 0x08000000
	meta_mod   = 0x10000000
	keypad_mod = 0x20000000
)

type (

	// Used for batching qml.Changed calls
	qmlDispatch struct{ value, field interface{} }

	// This allows us to trigger a qml.Changed on a specific
	// line in the view so that only it is
	// re-rendered by qml
	lineStruct struct {
		Text string
	}
)

var (
	t *qmlfrontend
)

func htmlcol(c render.Colour) string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

type (
)

func main() {
	// Need to lock the OS thread as OSX GUI requires GUI stuff to run in the main thread
	runtime.LockOSThread()

	log.AddFilter("file", log.FINEST, log.NewConsoleLogWriter())
	defer func() {
		py.NewLock()
		py.Finalize()
	}()

	t = &qmlfrontend{windows: make(map[*backend.Window]*frontendWindow)}
	go t.qmlBatchLoop()
	qml.Run(t.loop)
}
