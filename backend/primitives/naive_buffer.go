package primitives

const chunk_size = 256 * 1024

type (
	NaiveBuffer struct {
		data []rune
	}
)

func (b *NaiveBuffer) Size() int {
	return len(b.data)
}

func (buf *NaiveBuffer) Index(pos int) rune {
	return buf.data[pos]
}

func (buf *NaiveBuffer) SubstrR(r Region) []rune {
	l := len(buf.data)
	a, b := Clamp(0, l, r.Begin()), Clamp(0, l, r.End())
	return buf.data[a:b]
}

func (buf *NaiveBuffer) InsertR(point int, value []rune) {
	point = Clamp(0, len(buf.data), point)
	req := len(buf.data) + len(value)
	if cap(buf.data) < req {
		alloc := (req + chunk_size - 1) &^ (chunk_size - 1)
		n := make([]rune, len(buf.data), alloc)
		copy(n, buf.data)
		buf.data = n
	}
	if point == len(buf.data) {
		copy(buf.data[point:req], value)
	} else {
		copy(buf.data[point+len(value):cap(buf.data)], buf.data[point:len(buf.data)])
		copy(buf.data[point:req], value)
	}
	buf.data = buf.data[:req]
}

func (buf *NaiveBuffer) Erase(point, length int) {
	if length == 0 {
		return
	}
	buf.data = append(buf.data[0:point], buf.data[point+length:len(buf.data)]...)
}

func (b *NaiveBuffer) RowCol(point int) (row, col int) {
	if point < 0 {
		point = 0
	} else if l := b.Size(); point > l {
		point = l
	}

	sub := b.SubstrR(Region{0, point})
	for _, r := range sub {
		if r == '\n' {
			row++
			col = 0
		} else {
			col++
		}
	}
	return
}

func (b *NaiveBuffer) TextPoint(row, col int) (i int) {
	if row == 0 && col == 0 {
		return 0
	}
	for l := b.Size(); row > 0 && i < l; i++ {
		if b.data[i] == '\n' {
			row--
		}
	}
	if i < b.Size() {
		return i + col
	}
	return i
}
