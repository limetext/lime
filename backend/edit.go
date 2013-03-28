package backend

type (
	Action interface {
		Apply()
		Undo()
	}

	compositeAction struct {
		actions []Action
	}

	insertAction struct {
		buffer *Buffer
		point  int
		value  string
	}

	eraseAction struct {
		insertAction
		region Region
	}
)

func (ca *compositeAction) Apply() {
	for _, a := range ca.actions {
		a.Apply()
	}
}

func (ca *compositeAction) Undo() {
	l := len(ca.actions) - 1
	for i := range ca.actions {
		ca.actions[l-i].Undo()
	}
}

func (ia *insertAction) Apply() {
	ia.buffer.data = ia.buffer.data[0:ia.point] + ia.value + ia.buffer.data[ia.point:len(ia.buffer.data)]
}

func (ia *insertAction) Undo() {
	ia.buffer.data = ia.buffer.data[0:ia.point] + ia.buffer.data[ia.point+len(ia.value):len(ia.buffer.data)]
}

func (ea *eraseAction) Apply() {
	ea.region = ea.region.Clip(Region{0, ea.buffer.Size()})
	ea.value = ea.buffer.Substr(ea.region)
	ea.point = ea.region.Begin()
	ea.insertAction.Undo()
}

func (ea *eraseAction) Undo() {
	ea.insertAction.Apply()
}

func NewEraseAction(b *Buffer, region Region) Action {
	return &eraseAction{insertAction{buffer: b}, region}
}

func NewInsertAction(b *Buffer, point int, value string) Action {
	return &insertAction{b, clamp(0, len(b.data), point), value}
}

func NewReplaceAction(b *Buffer, region Region, value string) Action {
	return &compositeAction{[]Action{
		NewEraseAction(b, region),
		NewInsertAction(b, clamp(0, b.Size()-region.Size(), region.Begin()), value),
	}}
}
