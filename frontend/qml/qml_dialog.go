// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"github.com/limetext/lime/backend/log"
	"gopkg.in/qml.v1"
	"sync"
)

type (
	qmlDialog struct{}
)

func (q *qmlDialog) Show(msg, icon string) (ret int) {
	src := `import QtQuick 2.2
import QtQuick.Dialogs 1.1

Item {MessageDialog {
	objectName: "realDialog"
	id: messageDialog
	title: "May I have your attention please"
	text: "` + msg + `"
	icon: ` + icon + `
	standardButtons: StandardButton.Ok | StandardButton.Cancel
	Component.onCompleted: visible = true
}}`
	engine := qml.NewEngine()
	engine.Context().SetVar("q", q)
	component, err := engine.LoadString("dialog.qml", src)
	if err != nil {
		log.Errorf("Unable to instanciate dialog: %s", err)
		return 0
	}
	var wg sync.WaitGroup
	wg.Add(1)
	obj := component.Create(nil)
	obj = obj.ObjectByName("realDialog")
	obj.On("accepted", func() {
		ret = 1
		wg.Done()
	})
	obj.On("rejected", func() {
		ret = 0
		wg.Done()
	})

	wg.Wait()
	engine.Destroy()
	log.Debug("returning %d", ret)
	return
}
