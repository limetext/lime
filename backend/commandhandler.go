// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	. "github.com/limetext/lime/backend/util"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type (
	CommandHandler interface {
		Unregister(string) error
		Register(name string, cmd interface{}) error
		// TODO(q): Do the commands need to be split in separate lists?
		RunWindowCommand(*Window, string, Args) error
		RunTextCommand(*View, string, Args) error
		RunApplicationCommand(string, Args) error
	}

	appcmd         map[string]interface{}
	textcmd        map[string]interface{}
	wndcmd         map[string]interface{}
	commandHandler struct {
		ApplicationCommands appcmd
		TextCommands        textcmd
		WindowCommands      wndcmd
		log                 bool
		verbose             bool
	}
)

var casere = regexp.MustCompile(`([A-Z])`)

func PascalCaseToSnakeCase(in string) string {
	first := true
	return casere.ReplaceAllStringFunc(in, func(in string) string {
		if first {
			first = false
			return strings.ToLower(in)
		}
		return "_" + strings.ToLower(in)
	})

}

// If the cmd implements the CustomInit interface, its Init function
// is called, otherwise the fields of th cmd's underlying struct type
// will be enumerated and match against the dictionary keys in args,
// or if the key isn't provided in args, the Zero value will be used.
func (ch *commandHandler) init(cmd interface{}, args Args) error {
	if in, ok := cmd.(CustomInit); ok {
		return in.Init(args)
	}
	v := reflect.ValueOf(cmd).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		f := v.Field(i)
		if ft.Anonymous || !f.CanSet() {
			continue
		}
		key := PascalCaseToSnakeCase(ft.Name)
		fv, ok := args[key]
		if !ok {
			fv = reflect.Zero(ft.Type).Interface()
		}
		if f.CanAddr() {
			if f2, ok := f.Addr().Interface().(CustomSet); ok {
				if err := f2.Set(fv); err != nil {
					return err
				}
				continue
			}
		}
		f.Set(reflect.ValueOf(fv))
	}
	return nil
}

func (ch *commandHandler) RunWindowCommand(wnd *Window, name string, args Args) error {
	lvl := log4go.FINE
	p := Prof.Enter("wc")
	defer p.Exit()
	if ch.log {
		lvl = log4go.DEBUG
	}
	log4go.Logf(lvl, "Running window command: %s %v", name, args)
	t := time.Now()
	if c, ok := ch.WindowCommands[name].(WindowCommand); c != nil && ok {
		if err := ch.init(c, args); err != nil && ch.verbose {
			log4go.Debug("Command initialization failed: %s", err)
		} else if err := wnd.runCommand(c, name); err != nil {
			log4go.Logf(lvl+1, "Command execution failed: %s", err)
		} else {
			log4go.Logf(lvl, "Ran Window command: %s %s", name, time.Since(t))
		}
	} else {
		log4go.Logf(lvl, "No such window command: %s", name)
	}
	return nil
}

func (ch *commandHandler) RunTextCommand(view *View, name string, args Args) error {
	lvl := log4go.FINE
	p := Prof.Enter("tc")
	defer p.Exit()
	t := time.Now()
	if ch.log {
		lvl = log4go.DEBUG
	}
	log4go.Logf(lvl, "Running text command: %s %v", name, args)
	if c, ok := ch.TextCommands[name].(TextCommand); c != nil && ok {
		if err := ch.init(c, args); err != nil && ch.verbose {
			log4go.Debug("Command initialization failed: %s", err)
		} else if err := view.runCommand(c, name); err != nil {
			log4go.Logf(lvl, "Command execution failed: %s", err)
		}
	} else if w := view.Window(); w != nil {
		if c, ok := ch.WindowCommands[name].(WindowCommand); c != nil && ok {
			if err := w.runCommand(c, name); err != nil {
				log4go.Logf(lvl, "Command execution failed: %s", err)
			}
		}
	}
	log4go.Logf(lvl, "Ran text command: %s %s", name, time.Since(t))
	return nil
}

func (ch *commandHandler) RunApplicationCommand(name string, args Args) error {
	p := Prof.Enter("ac")
	defer p.Exit()
	if ch.log {
		log4go.Info("Running application command: %s %v", name, args)
	} else {
		log4go.Fine("Running application command: %s %v", name, args)
	}
	if c, ok := ch.ApplicationCommands[name].(ApplicationCommand); c != nil && ok {
		if err := ch.init(c, args); err != nil && ch.verbose {
			log4go.Debug("Command initialization failed: %s", err)
		} else if err := c.Run(); err != nil && ch.verbose {
			log4go.Debug("Command execution failed: %s", err)
		}
	}
	return nil
}

func (ch *commandHandler) Unregister(name string) error {
	if _, ok := ch.ApplicationCommands[name]; ok {
		ch.ApplicationCommands[name] = nil
	} else if _, ok := ch.TextCommands[name]; ok {
		ch.TextCommands[name] = nil
	} else if _, ok := ch.WindowCommands[name]; !ok {
		return fmt.Errorf("%s wasn't a registered command", name)
	} else {
		ch.WindowCommands[name] = nil
	}
	return nil
}

func (ch *commandHandler) Register(name string, cmd interface{}) error {
	var r = false
	log4go.Finest("Want to register %s", name)
	if ac, ok := cmd.(ApplicationCommand); ok {
		if _, ok := ch.ApplicationCommands[name]; ok {
			return fmt.Errorf("%s is already a registered command", name)
		}
		r = true
		ch.ApplicationCommands[name] = ac
	}
	if wc, ok := cmd.(WindowCommand); ok {
		if _, ok := ch.WindowCommands[name]; ok {
			return fmt.Errorf("%s is already a registered command", name)
		}
		r = true
		ch.WindowCommands[name] = wc
	}
	if tc, ok := cmd.(TextCommand); ok {
		if _, ok := ch.TextCommands[name]; ok {
			return fmt.Errorf("%s is already a registered command", name)
		}
		r = true
		ch.TextCommands[name] = tc
	}
	if !r {
		return fmt.Errorf("Command wasn't registered in any list: %s", name)
	} else if ch.verbose {
		log4go.Finest("Successfully registered command %s", name)
	}
	return nil
}
