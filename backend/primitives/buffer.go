package primitives

import (
	"regexp"
	"strings"
)

type (
	Buffer struct {
		HasId
		changecount int
		name        string
		filename    string
		data        string
		callbacks   []BufferChangedCallback
	}
	BufferChangedCallback func(buf *Buffer, position, delta int)
)

func (b *Buffer) AddCallback(cb BufferChangedCallback) {
	b.callbacks = append(b.callbacks, cb)
}

func (b *Buffer) SetName(n string) {
	b.name = n
}

func (b *Buffer) Name() string {
	return b.name
}

func (b *Buffer) FileName() string {
	return b.filename
}

func (b *Buffer) SetFileName(n string) {
	b.filename = n
}

func (b *Buffer) Size() int {
	return len(b.data)
}

func (buf *Buffer) Substr(r Region) string {
	l := len(buf.data)
	a, b := Clamp(0, l, r.Begin()), Clamp(0, l, r.End())
	return string(buf.data[a:b])
}

func (buf *Buffer) notify(position, delta int) {
	for i := range buf.callbacks {
		buf.callbacks[i](buf, position, delta)
	}
}

func (buf *Buffer) Insert(point int, value string) {
	if len(value) == 0 {
		return
	}
	buf.changecount++
	buf.data = buf.data[0:point] + value + buf.data[point:len(buf.data)]
	buf.notify(point, len(value))
}

func (buf *Buffer) Erase(point, length int) {
	if length == 0 {
		return
	}
	buf.changecount++
	buf.data = buf.data[0:point] + buf.data[point+length:len(buf.data)]
	buf.notify(point+length, -length)
}

func (b *Buffer) Data() string {
	return b.data
}

func (b *Buffer) ChangeCount() int {
	return b.changecount
}

func (b *Buffer) RowCol(point int) (row, col int) {
	if point < 0 {
		point = 0
	} else if l := b.Size(); point > l {
		point = l
	}
	lines := strings.Split(b.data[:point], "\n")
	if l := len(lines); l == 0 {
		return 1, 1
	} else {
		return l, len(lines[l-1]) + 1
	}
}

func (b *Buffer) TextPoint(row, col int) int {
	lines := strings.Split(b.data, "\n")
	if row < 1 || len(lines) == 0 {
		return 0
	}
	if col == 0 {
		col = 1
	}
	if row == 1 {
		col -= 1
	}
	if row > len(lines) {
		return b.Size()
	}
	offset := len(strings.Join(lines[:row-1], "\n")) + col
	return offset
}

func (b *Buffer) Line(offset int) Region {
	if offset < 0 {
		return Region{0, 0}
	} else if s := b.Size(); offset >= s {
		return Region{s, s}
	} else if b.data[offset] == '\n' {
		return Region{offset, offset}
	}
	data := b.data
	s := offset
	for s > 0 && data[s-1] != '\n' {
		s--
	}
	e := offset + 1
	for e < len(data) && data[e] != '\n' {
		e++
	}
	return Region{s, e}
}

func (b *Buffer) FullLine(offset int) Region {
	r := b.Line(offset)
	d := b.data
	s := b.Size()
	for r.B < s && (d[r.B] == '\r' || d[r.B] == '\n') {
		r.B++
	}
	return r
}

var (
	vwre1 = regexp.MustCompile(`\b\w*$`)
	vwre2 = regexp.MustCompile(`^\w*`)
)

func (b *Buffer) Word(offset int) Region {
	_, col := b.RowCol(offset)
	lr := b.Line(offset)
	line := b.Substr(lr)
	begin := 0
	end := len(line)

	if col > len(line) {
		col = len(line)
	}
	if m := vwre1.FindStringIndex(line[:col]); m != nil {
		begin = m[0]
	} else {
		return Region{offset, offset}
	}
	if m := vwre2.FindStringIndex(line[begin:]); m != nil {
		end = begin + m[1]
	}
	return Region{lr.Begin() + begin, lr.Begin() + end}
}
