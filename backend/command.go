package backend

type (
	// The Args type is just a generic string key and interface{} value
	// map type used to serialize command arguments.
	Args map[string]interface{}

	// The CustomSet interface can be optionally implemented
	// by struct members of a concrete command struct.
	//
	// If implemented, it'll be called by the default
	// Command initialization code with the data gotten
	// from the Args map.
	CustomSet interface {
		Set(v interface{}) error
	}

	// The CustomInit interface can be optionally implemented
	// by a Command and will be called instead of the default
	// command initialization code.
	CustomInit interface {
		Init(args Args) error
	}

	// The Command interface implements the basic interface
	// that is shared between the different more specific
	// command type interfaces.
	Command interface {
		// Returns whether the Command is enabled or not.
		IsEnabled() bool

		// Returns whether the Command is visible in menus,
		// the goto anything panel or other user interface
		// that lists available commands.
		IsVisible() bool

		// Returns the textual description of the command.
		Description() string

		// Whether or not this Command bypasses the undo stack.
		BypassUndo() bool
	}

	// The WindowCommand interface extends the base Command interface
	// with functionality specific for WindowCommands.
	WindowCommand interface {
		Command

		// Execute this command with the specified window as the
		// argument
		Run(*Window) error
	}

	// The TextCommand interface extends the base Command interface
	// with functionality specific for TextCommands.
	TextCommand interface {
		Command
		// Execute this command with the specified View and Edit object
		// as the arguments
		Run(*View, *Edit) error
	}

	// The ApplicationCommand interface extends the base Command interface
	// with functionality specific for ApplicationCommands.
	ApplicationCommand interface {
		Command
		// Execute this command
		Run() error
		// Returns whether this command is checked or not.
		// Used to display a checkbox in the user interface
		// for boolean commands.
		IsChecked() bool
	}

	// The DefaultCommand implements the default operation
	// of the basic Command interface and is recommended to
	// be used as the base when creating new Commands.
	DefaultCommand struct{}

	// The BypassUndoCommand is the same as the DefaultCommand
	// type, except that its implementation of BypassUndo returns
	// true rather than false.
	BypassUndoCommand struct {
		DefaultCommand
	}
)

// The default is to not bypass the undo stack.
func (d *DefaultCommand) BypassUndo() bool {
	return false
}

// By default a command is enabled.
func (d *DefaultCommand) IsEnabled() bool {
	return true
}

// By default a command is visible.
func (d *DefaultCommand) IsVisible() bool {
	return true
}

// By default the string "TODO" is return as the description.
func (d *DefaultCommand) Description() string {
	return "TODO"
}

// The BypassUndoCommand defaults to bypassing the
// undo stack.
func (b *BypassUndoCommand) BypassUndo() bool {
	return true
}
