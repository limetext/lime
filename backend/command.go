package backend

type (
	Args    map[string]interface{}
	Command interface {
		IsEnabled(Args) bool
		IsVisible(Args) bool
		Description(Args) string
	}
	WindowCommand interface {
		Command
		Run(Args)
	}
	TextCommand interface {
		Command
		Run(Edit, Args)
	}
	ApplicationCommand interface {
		Command
		Run(Args)
		IsChecked(Args) bool
	}

	DefaultCommand struct{}
)

func (d *DefaultCommand) IsEnabled(Args) bool {
	return true
}

func (d *DefaultCommand) IsVisible(Args) bool {
	return true
}

func (d *DefaultCommand) Description(Args) string {
	return "TODO"
}
