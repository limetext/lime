package backend

type (
	Args    map[string]interface{}
	Command interface {
		IsEnabled(Args) bool
		IsVisible(Args) bool
		Description(Args) string
		BypassUndo() bool
	}
	WindowCommand interface {
		Command
		Run(*Window, Args) error
	}
	TextCommand interface {
		Command
		Run(*View, *Edit, Args) error
	}
	ApplicationCommand interface {
		Command
		Run(Args) error
		IsChecked(Args) bool
	}

	DefaultCommand struct{}
	// Use with caution..
	BypassUndoCommand struct {
		DefaultCommand
	}
)

func (d *DefaultCommand) BypassUndo() bool {
	return false
}

func (d *DefaultCommand) IsEnabled(Args) bool {
	return true
}

func (d *DefaultCommand) IsVisible(Args) bool {
	return true
}

func (d *DefaultCommand) Description(Args) string {
	return "TODO"
}

func (b *BypassUndoCommand) BypassUndo() bool {
	return true
}
