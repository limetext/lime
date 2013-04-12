package primitives

import (
	//	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

var complexnode_test = &node{
	22,
	&node{
		9,
		&node{
			9,
			&node{
				6,
				&node{6, nil, nil, []rune("Hello ")},
				&node{3, nil, nil, []rune("my ")},
				nil,
			},
			nil,
			nil,
		},
		&node{
			7,
			&node{
				6,
				&node{
					2,
					&node{2, nil, nil, []rune("na")},
					&node{4, nil, nil, []rune("me i")},
					nil,
				},
				&node{1, nil, nil, []rune("s")},
				nil,
			},
			&node{6, nil, nil, []rune(" Simon")},
			nil,
		},
		nil,
	},
	nil,
	nil,
}

type Test struct {
	in  *node
	exp string
}

var (
	tests = []Test{
		{&node{6, &node{6, nil, nil, []rune("Hello ")}, &node{5, nil, nil, []rune("world")}, nil}, "Hello world"},
		{&node{6, &node{6, nil, nil, []rune("Hello ")}, &node{3, &node{3, nil, nil, []rune("wor")}, &node{2, nil, nil, []rune("ld")}, nil}, nil}, "Hello world"},
		{&node{6, &node{6, nil, nil, []rune("Hello ")}, &node{5, nil, nil, []rune("world")}, nil}, "Hello world"},
		{complexnode_test, "Hello my name is Simon"},
	}
	merges = []int{4, 8, 32, 128, 1024, merge}
)

func init() {
	const (
		size  = 1024
		split = 8
	)
	in := make([]rune, size)
	fill(in)
	tests = append(tests, Test{newNodeEx(in, split), string(in)})
}

func TestNode(t *testing.T) {
	for i, test := range tests {
		if sub := test.in.Substr(Region{0, len(test.exp)}); sub != test.exp {
			t.Fatalf("%d %s != %s", i, sub, test.exp)
		} else if l := test.in.length(); l != len(sub) {
			t.Fatalf("%d %d != %d", l, len(sub))
		}
	}
}

func TestNodeSimplify(t *testing.T) {
	r := &node{5, nil, nil, []rune("world")}
	l := &node{0, nil, nil, nil}
	n := node{0, l, r, nil}
	n.simplify()
	if !reflect.DeepEqual(&n, r) {
		t.Error(n.dump(""))
	}
	n = node{5, r, l, nil}
	n.simplify()
	if !reflect.DeepEqual(&n, r) {
		t.Error(n.dump(""))
	}
}

func TestNodeSplit(t *testing.T) {
	for _, m := range merges {
		merge = m
		for i, test := range tests {
			for j := range test.exp {
				nn := test.in.clone()
				r := nn.split(j)
				if sub := nn.Substr(Region{0, j}) + r.Substr(Region{0, len(test.exp) - j}); sub != test.exp {
					t.Fatalf("%d, %d, split %s != %s:\n%s\n%s", i, j, sub, test.exp, nn.dump("\t"), r.dump("\t"))
				} else if l := nn.length(); l != j {
					t.Fatalf("%d, %d, split length1 %d != %d:\n%s\n%s", i, j, l, j, nn.dump("\t"), r.dump("\t"))
				} else if l := r.length(); l != len(test.exp)-j {
					t.Fatalf("%d, %d, split length2 %d != %d:\n%s\n", i, j, l, len(test.exp)-j, r.dump("\t"))
				}
			}
		}
	}
}

func TestNodeConcat(t *testing.T) {
	for _, m := range merges {
		merge = m
		for i, test := range tests {
			for j := range test.exp {
				nn := test.in.clone()
				r := nn.split(j)
				nn.concat(r)
				if sub := nn.Substr(Region{0, len(test.exp)}); sub != test.exp {
					t.Fatalf("%d, %d, split/concat %s != %s:\n%s", i, j, sub, test.exp, nn.dump("\t"))
				} else if l := nn.length(); l != len(test.exp) {
					t.Fatalf("%d, %d, %d != %d", i, j, l, len(test.exp))
				}
			}
		}
	}
}

func TestNodeInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("Short")
	}

	const (
		size  = 256
		isize = 256
	)
	od := make([]rune, size)
	fill(od)

	in := make([]rune, size)
	fill(in)

	ins := string(in)

	for _, m := range merges {
		merge = m
		for i := range od {
			n := newNode(od)
			od2 := make([]rune, size)
			copy(od2, od)
			b := Buffer{data: od2}
			n.insert(i, ins)
			b.Insert(i, ins)
			r := Region{0, b.Size()}

			if b.Size() != n.length() {
				na := n.dump("\t")
				t.Fatalf("%d, %d: %d != %d\n%s", m, i, b.Size(), n.length(), na)
			} else if e, a := b.Substr(r), n.Substr(r); e != a {
				na := n.dump("\t")
				t.Fatalf("%d, %d: %s != %s\n%s", m, i, e, a, na)
			}
		}
		for i := range od {
			n := newNode(od)
			for _, j := range ins {
				l := n.length()
				n.insert(i, string(j))
				if n.length() != l+1 {
					t.Log(string(j))
					na := n.dump("\t")
					t.Fatalf("%d, %d, %d: %d != %d\n%s", m, i, j, n.length(), l+1, na)
				}
			}
		}
	}
}

func TestNodeErase(t *testing.T) {
	if testing.Short() {
		t.Skip("Short")
	}

	const (
		size  = 2 * 1024
		dsize = 1
	)
	od := make([]rune, size)
	fill(od)

	for _, m := range merges {
		merge = m
		for i := range od {
			n := newNode(od)
			od2 := make([]rune, size)
			copy(od2, od)
			b := Buffer{data: od2}
			n.erase(i, dsize)
			b.Erase(i, dsize)
			r := Region{0, b.Size()}

			if b.Size() != n.length() {
				t.Fatalf("%d, %d: %d != %d\n%s", m, i, b.Size(), n.length())
			} else if e, a := b.Substr(r), n.Substr(r); e != a {
				r = Region{0, 20}
				e = b.Substr(r)
				a = n.Substr(r)
				t.Fatalf("%d, %d: %s != %s (%v)", m, i, e, a, e != a)
			}
		}
	}
}

func BenchmarkNodeSplit(b *testing.B) {
	b.StopTimer()
	data := make([]rune, 1024*256)
	fill(data)
	buf := newNodeEx(data, 4096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b2 := buf.clone()
		b2.split(512)
	}
}

var _ = rand.ExpFloat64

func BenchmarkNodeInsertRand(b *testing.B) {
	r := rand.Perm(b.N)
	b.StopTimer()
	sdata := testinsert()
	buf := newNode(testbuffer().Runes())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l := buf.length()
		pos := r[i] % l
		buf.insert(pos, sdata)
	}
	buf.Substr(Region{0, buf.length()})
}

func BenchmarkNodeInsertBegin(b *testing.B) {
	b.StopTimer()
	sdata := testinsert()
	in := testbuffer().Runes()
	buf := newNode(in)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.insert(0, sdata)
	}
	buf.Substr(Region{0, buf.length()})
	if a, e := buf.length(), b.N*len(sdata)+len(in); a != e {
		b.Error(a, e)
	}
}

func BenchmarkNodeInsertMid(b *testing.B) {
	b.StopTimer()
	sdata := testinsert()
	in := testbuffer().Runes()
	buf := newNode(in)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.insert(buf.length()/2, sdata)
	}
	buf.Substr(Region{0, buf.length()})
	if a, e := buf.length(), b.N*len(sdata)+len(in); a != e {
		b.Error(a, e)
	}
}

func BenchmarkNodeInsertEnd(b *testing.B) {
	b.StopTimer()
	sdata := testinsert()
	in := testbuffer().Runes()
	buf := newNode(in)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.insert(buf.length(), sdata)
	}
	buf.Substr(Region{0, buf.length()})
	if a, e := buf.length(), b.N*len(sdata)+len(in); a != e {
		b.Error(a, e)
	}
}
