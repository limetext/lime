package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	. "lime/backend/util"
	"time"
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
		verbose             bool
	}
)

func (ch *commandHandler) RunWindowCommand(wnd *Window, name string, args Args) error {
	lvl := log4go.FINE
	p := Prof.Enter("wc")
	defer p.Exit()
	if ch.log {
		lvl = log4go.DEBUG
	}
	log4go.Logf(lvl, "Running window command: %s %v", name, args)
	t := time.Now()
	if c := ch.WindowCommands[name]; c != nil {
		if err := wnd.runCommand(c, name, args); err != nil {
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
	if c := ch.TextCommands[name]; c != nil {
		if err := view.runCommand(c, name, args); err != nil {
			log4go.Logf(lvl, "Command execution failed: %s", err)
		}
	} else if w := view.Window(); w != nil {
		if c := ch.WindowCommands[name]; c != nil {
			if err := w.runCommand(c, name, args); err != nil {
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
	if c := ch.ApplicationCommands[name]; c != nil {
		if err := c.Run(args); err != nil && ch.verbose {
			log4go.Debug("Command execution failed: %s", err)
		}
	}
	return nil
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
