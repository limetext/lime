package backend

type (
	Editor interface {
		Windows() []*Window
		ActiveWindow() *Window
		Arch() string
		Platform() string
		Version() string
		LogInput(bool)
		LogCommands(bool)
		GetClipboard() string
		SetClipboard(string)
	}
)
