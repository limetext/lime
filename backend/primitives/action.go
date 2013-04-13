package primitives

import (
	"fmt"
)

type (
	Action interface {
		Apply()
		Undo()
	}

	CompositeAction struct {
		actions []Action
	}

	insertAction struct {
		buffer *Buffer
		point  int
		value  []rune
	}

	eraseAction struct {
		insertAction
		region Region
	}
)

func (ca CompositeAction) String() string {
	ret := fmt.Sprintf("%d actions:\n", len(ca.actions))
	for i := range ca.actions {
		ret += fmt.Sprintf("\t%s\n", ca.actions[i])
	}
	return ret
}

func (ca *CompositeAction) Apply() {
	for _, a := range ca.actions {
		a.Apply()
	}
}

func (ca *CompositeAction) Undo() {
	l := len(ca.actions) - 1
	for i := range ca.actions {
		ca.actions[l-i].Undo()
	}
}

func (ca *CompositeAction) Add(a Action) {
	ca.actions = append(ca.actions, a)
}

func (ca *CompositeAction) AddExec(a Action) {
	ca.Add(a)
	ca.actions[len(ca.actions)-1].Apply()
}

func (ca *CompositeAction) Len() int {
	return len(ca.actions)
}

func (ia *insertAction) Apply() {
	ia.buffer.Insert(ia.point, string(ia.value))
}

func (ia *insertAction) Undo() {
	ia.buffer.Erase(ia.point, len(ia.value))
}

func (ea *eraseAction) Apply() {
	ea.region = ea.region.Clip(Region{0, ea.buffer.Size()})
	ea.value = []rune(ea.buffer.Substr(ea.region))
	ea.point = ea.region.Begin()
	ea.insertAction.Undo()
}

func (ea *eraseAction) Undo() {
	ea.insertAction.Apply()
}

func (ia insertAction) String() string {
	return fmt.Sprintf("insert %d %s", ia.point, string(ia.value))
}

func (ea eraseAction) String() string {
	return fmt.Sprintf("erase %v", ea.region)
}

func NewEraseAction(b *Buffer, region Region) Action {
	return &eraseAction{insertAction{buffer: b}, region}
}

func NewInsertAction(b *Buffer, point int, value string) Action {
	return &insertAction{b, Clamp(0, b.Size(), point), []rune(value)}
}

func NewReplaceAction(b *Buffer, region Region, value string) Action {
	return &CompositeAction{[]Action{
		NewEraseAction(b, region),
		NewInsertAction(b, Clamp(0, b.Size()-region.Size(), region.Begin()), value),
	}}
}
