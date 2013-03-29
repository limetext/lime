package backend

type (
	KeyPress struct {
		KeyCode, RawKeyCode     int
		Shift, Super, Alt, Ctrl bool
	}

	KeyContext struct {
		Key, Operator string
		Operand       interface{}
		Match_all     bool
	}

	KeyBinding struct {
		Keys    []KeyPress
		Command string
		Args    map[string]interface{}
		Context []KeyContext
	}
)
