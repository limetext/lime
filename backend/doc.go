// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// The backend package defines the very soul of Lime.
//
// Some highlevel concepts follow.
//
// Frontend
//
// Lime is designed with the goal of having a clear frontend
// and backend separation to allow and hopefully simplify
// the creation of multiple frontend versions.
//
// The two most active frontends are at the time of writing this
// one for terminals based on termbox-go and a GUI application
// based on Qt's QML scripting language.
//
// There's also a Proof of Concept html frontend.
//
// Editor
//
// The Editor singleton represents the most fundamental interface
// that frontends use to access the backend. It keeps a list of
// editor windows, handles input, detects file changes as well
// as communicate back to the frontend as needed.
//
// Window
//
// At any time there can be multiple windows of the editor open.
// Each Window can have a different layout, settings and status.
//
// View
//
// The View class defines a "view" into a specific backing buffer.
// Multiple views can share the same backing buffer. For instance
// viewing the same buffer in split view, or viewing the buffer
// with one syntax definition in one view and another syntax
// definition in the other.
//
// It is very closely related to the view defined in
// the Model-view-controller paradigm, and contains settings
// pertaining to exactly how a buffer is shown to the user.
//
// Command
//
// The command interface defines actions to be executed either
// for the whole application, a specific window or a specific
// view.
//
// Key bindings
//
// Key bindings define a sequence of key-presses, a Command and
// the command's arguments to be executed upon that sequence having
// been pressed.
//
// Key bindings can optionally have multiple contexts associated with it
// which allows the exact same key sequence to have different meaning
// depending on context.
//
// See http://godoc.org/github.com/limetext/lime/backend#QueryContextCallback
// for details.
//
// Settings
//
// Many of the components have their own key-value Settings object associated with
// it, but settings are also nested. In other words, if the settings key does not
// exist in the current object's settings, its parent's settings object is queried
// next which in turn will query its parent if its settings object didn't contain the
// key neither.
//
package backend
