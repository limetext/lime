package backend

import (
	"code.google.com/p/log4go"
	"fmt"
)

type (
	CommandHandler interface {
		Unregister(string) error
		Register(name string, cmd Command) error
		// TODO(q): Do the commands need to be split in separate lists?
		RunWindowCommand(*Window, string, Args) error
		RunTextCommand(*View, string, Args) error
		RunApplicationCommand(string, Args) error
	}

	appcmd         map[string]ApplicationCommand
	textcmd        map[string]TextCommand
	wndcmd         map[string]WindowCommand
	commandHandler struct {
		ApplicationCommands appcmd
		TextCommands        textcmd
		WindowCommands      wndcmd
		log                 bool
	}
)

func (ch *commandHandler) RunWindowCommand(wnd *Window, name string, args Args) error {
	if ch.log {
		log4go.Debug("RunWindowCommand: %s %v", name, args)
	}
	if wc, ok := ch.WindowCommands[name]; ok {
		return wc.Run(wnd, args)
	} else {
		return fmt.Errorf("No such command registered: %s", name)
	}
}

func (ch *commandHandler) RunTextCommand(view *View, name string, args Args) error {
	if ch.log {
		log4go.Debug("RunTextCommand: %s %v", name, args)
	}

	if tc, ok := ch.TextCommands[name]; ok {
		e := view.BeginEdit()
		err := tc.Run(view, e, args)
		view.EndEdit(e)
		return err
	} else {
		return fmt.Errorf("No such command registered: %s", name)
	}
}

func (ch *commandHandler) RunApplicationCommand(name string, args Args) error {
	if ch.log {
		log4go.Debug("RunTextCommand: %s %v", name, args)
	}

	if ac, ok := ch.ApplicationCommands[name]; ok {
		return ac.Run(args)
	} else {
		return fmt.Errorf("No such command registered: %s", name)
	}
}

func (ch *commandHandler) Unregister(name string) error {
	if _, ok := ch.ApplicationCommands[name]; !ok {
		if _, ok := ch.TextCommands[name]; !ok {
			if _, ok := ch.WindowCommands[name]; !ok {
				return fmt.Errorf("%s wasn't a registered command", name)
			} else {
				ch.WindowCommands[name] = nil
			}
		} else {
			ch.TextCommands[name] = nil
		}
	} else {
		ch.ApplicationCommands[name] = nil
	}
	return nil
}

func (ch *commandHandler) Register(name string, cmd Command) error {
	if ac, ok := cmd.(ApplicationCommand); ok {
		if _, ok := ch.ApplicationCommands[name]; ok {
			return fmt.Errorf("%s is already a registered command", name)
		}
		ch.ApplicationCommands[name] = ac
	}
	if wc, ok := cmd.(WindowCommand); ok {
		if _, ok := ch.WindowCommands[name]; ok {
			return fmt.Errorf("%s is already a registered command", name)
		}
		ch.WindowCommands[name] = wc
	}
	if tc, ok := cmd.(TextCommand); ok {
		if _, ok := ch.TextCommands[name]; ok {
			return fmt.Errorf("%s is already a registered command", name)
		}
		ch.TextCommands[name] = tc
	}
	return nil
}
