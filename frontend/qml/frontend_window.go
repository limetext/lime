// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"github.com/limetext/lime/backend"
	"gopkg.in/qml.v1"
	"sync"
)

// A helper glue structure connecting the backend Window with the qml.Window
type frontendWindow struct {
	bw     *backend.Window
	views  []*frontendView
	window *qml.Window
}

// Instantiates a new window, and launches a new goroutine waiting for it
// to be closed. The WaitGroup is increased at function entry and decreased
// once the window closes.
func (fw *frontendWindow) launch(wg *sync.WaitGroup, component qml.Object) {
	wg.Add(1)
	fw.window = component.CreateWindow(nil)
	fw.window.Show()
	fw.window.Set("myWindow", fw)

	go func() {
		fw.window.Wait()
		wg.Done()
	}()
}

func (fw *frontendWindow) View(idx int) *frontendView {
	return fw.views[idx]
}

func (fw *frontendWindow) ActiveViewIndex() int {
	for i, v := range fw.views {
		if v.bv == fw.bw.ActiveView() {
			return i
		}
	}
	return 0
}

func (fw *frontendWindow) Back() *backend.Window {
	return fw.bw
}
