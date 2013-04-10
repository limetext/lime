package primitives

import (
	"regexp"
	"strings"
)

type (
	Buffer struct {
		HasId
		HasSettings
		changecount int
		name        string
		filename    string
		data        []rune
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

func (buf *Buffer) Insert(point int, svalue string) {
	if len(svalue) == 0 {
		return
	}
	value := []rune(svalue)
	req := len(buf.data) + len(value)
	if cap(buf.data) < req {
		n := make([]rune, len(buf.data), req)
		copy(n, buf.data)
		buf.data = n
	}
	copy(buf.data[point+len(value):cap(buf.data)], buf.data[point:len(buf.data)])
	copy(buf.data[point:cap(buf.data)], value)
	buf.data = buf.data[:req]
	buf.changecount++
	buf.notify(point, len(value))
}

func (buf *Buffer) Erase(point, length int) {
	if length == 0 {
		return
	}
	buf.changecount++
	buf.data = append(buf.data[0:point], buf.data[point+length:len(buf.data)]...)
	buf.notify(point+length, -length)
}

func (b *Buffer) String() string {
	return string(b.data)
}

func (b *Buffer) Runes() []rune {
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
	lines := strings.Split(string(b.data[:point]), "\n")
	if l := len(lines); l == 0 {
		return 0, 0
	} else {
		return l - 1, len([]rune(lines[l-1]))
	}
}

func (b *Buffer) TextPoint(row, col int) int {
	lines := strings.Split(string(b.data), "\n")
	if row < 0 || len(lines) == 0 {
		return 0
	}
	if row > len(lines) {
		return b.Size()
	}
	if row == 0 {
		col--
	}
	offset := len([]rune(strings.Join(lines[:row], "\n"))) + col + 1
	return offset
}

func (b *Buffer) Line(offset int) Region {
	if offset < 0 {
		return Region{0, 0}
	} else if s := b.Size(); offset >= s {
		return Region{s, s}
	}
	data := b.data
	s := offset
	for s > 0 && data[s-1] != '\n' {
		s--
	}
	e := offset
	for e < len(data) && data[e] != '\n' {
		e++
	}
	return Region{s, e}
}

// Returns a region that starts at the first character in a line
// and ends with the last character in a (possibly different) line
func (b *Buffer) Lines(r Region) Region {
	s := b.Line(r.Begin())
	e := b.Line(r.End())
	return Region{s.Begin(), e.End()}
}

func (b *Buffer) FullLine(offset int) Region {
	r := b.Line(offset)
	d := b.data
	s := b.Size()
	for r.B < s && (d[r.B] != '\r' && d[r.B] != '\n') {
		r.B++
	}
	if r.B != b.Size() {
		r.B++
	}
	return r
}

// Returns a region that starts at the first character in a line
// and ends with the line break in a (possibly different) line
func (b *Buffer) FullLines(r Region) Region {
	s := b.FullLine(r.Begin())
	e := b.FullLine(r.End())
	return Region{s.Begin(), e.End()}
}

var (
	vwre1 = regexp.MustCompile(`\b\w*$`)
	vwre2 = regexp.MustCompile(`^\w*`)
)

func (b *Buffer) Word(offset int) Region {
	_, col := b.RowCol(offset)
	lr := b.FullLine(offset)

	line := b.data[lr.Begin():lr.End()]
	if len(line) == 0 {
		return Region{offset, offset}
	}

	seps := "./\\()\"'-:,.;<>~!@#$%^&*|+=[]{}`~?"
	if v, ok := b.Settings().Get("word_separators", seps).(string); ok {
		seps = v
	}
	spacing := " \n\t\r"
	eseps := seps + spacing

	if col >= len(line) {
		col = len(line) - 1
	}
	last := true
	li := 0
	ls := false
	lc := 0
	for i, r := range line {
		cur := strings.ContainsRune(eseps, r)
		cs := r == ' '
		if !cs {
			lc = i
		}
		if last != cur || ls != cs {
			ls = cs
			r := Region{li, i}
			if r.Contains(col) && i != 0 {
				r.A, r.B = r.A+lr.Begin(), r.B+lr.Begin()
				if !(r.B == offset && last) {
					return r
				}
			}
			li = i
			last = cur
		}
	}
	r := Region{lr.Begin() + li, lr.End()}
	lc += lr.Begin()
	if lc != offset && !strings.ContainsRune(spacing, b.data[r.A]) {
		r.B = lc
	}
	if r.A == offset && r.B == r.A+1 {
		r.B--
	}
	return r
}

// Returns a region that starts at the first character in a word
// and ends with the last character in a (possibly different) word
func (b *Buffer) Words(r Region) Region {
	s := b.Word(r.Begin())
	e := b.Word(r.End())
	return Region{s.Begin(), e.End()}
}
