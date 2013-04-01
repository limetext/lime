package backend

import (
	"code.google.com/p/log4go"
	"fmt"
)

type (
	CommandHandler interface {
		UnregisterApplicationCommand(string) error
		RegisterApplicationCommand(name string, cmd ApplicationCommand) error
		UnregisterWindowCommand(string) error
		RegisterWindowCommand(name string, cmd WindowCommand) error
		UnregisterTextCommand(string) error
		RegisterTextCommand(name string, cmd TextCommand) error
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

func (ch *commandHandler) UnregisterApplicationCommand(name string) error {
	if _, ok := ch.ApplicationCommands[name]; !ok {
		return fmt.Errorf("%s wasn't a registered command", name)
	}
	ch.ApplicationCommands[name] = nil
	return nil
}

func (ch *commandHandler) RegisterApplicationCommand(name string, cmd ApplicationCommand) error {
	if _, ok := ch.ApplicationCommands[name]; !ok {
		return fmt.Errorf("%s is already a registered command", name)
	}
	ch.ApplicationCommands[name] = nil
	return nil
}

func (ch *commandHandler) UnregisterWindowCommand(name string) error {
	if _, ok := ch.WindowCommands[name]; !ok {
		return fmt.Errorf("%s wasn't a registered command", name)
	}
	ch.WindowCommands[name] = nil
	return nil
}

func (ch *commandHandler) RegisterWindowCommand(name string, cmd WindowCommand) error {
	if _, ok := ch.WindowCommands[name]; !ok {
		return fmt.Errorf("%s is already a registered command", name)
	}
	ch.WindowCommands[name] = nil
	return nil
}

func (ch *commandHandler) UnregisterTextCommand(name string) error {
	if _, ok := ch.TextCommands[name]; !ok {
		return fmt.Errorf("%s wasn't a registered command", name)
	}
	ch.TextCommands[name] = nil
	return nil
}

func (ch *commandHandler) RegisterTextCommand(name string, cmd TextCommand) error {
	if _, ok := ch.TextCommands[name]; !ok {
		return fmt.Errorf("%s is already a registered command", name)
	}
	ch.TextCommands[name] = nil
	return nil
}
