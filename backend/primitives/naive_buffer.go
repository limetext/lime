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

func (buf *NaiveBuffer) Insert(point int, value []rune) {
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
