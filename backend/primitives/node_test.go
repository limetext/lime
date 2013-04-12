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

// func init() {
// 	merge = 6
// }

type Test struct {
	in  node
	exp string
}

var tests = []Test{
	{node{6, &node{6, nil, nil, []rune("Hello ")}, &node{5, nil, nil, []rune("world")}, nil}, "Hello world"},
	{node{6, &node{6, nil, nil, []rune("Hello ")}, &node{3, &node{3, nil, nil, []rune("wor")}, &node{2, nil, nil, []rune("ld")}, nil}, nil}, "Hello world"},
	{node{6, &node{6, nil, nil, []rune("Hello ")}, &node{5, nil, nil, []rune("world")}, nil}, "Hello world"},
	{*complexnode_test, "Hello my name is Simon"},
}

func TestNode(t *testing.T) {
	for _, test := range tests {
		if sub := test.in.Substr(Region{0, len(test.exp)}); sub != test.exp {
			t.Fatalf("%s != %s", sub, test.exp)
		} else if l := test.in.length(); l != len(sub) {
			t.Fatalf("%d != %d", l, len(sub))
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

func TestNodeConcat(t *testing.T) {
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

func TestNodeInsert(t *testing.T) {
	// Todo...
	n := complexnode_test.clone()
	n.insert(12, "krankelibrankelfnatt")
	t.Log(n, n.dump(""))
	n.erase(12, len("krankelibrankelfnatt"))
	t.Log(n, n.dump(""))

	t.Log(n, n.dump(""))
	for _, c := range "krankelibrankelfnatt" {
		n.insert(12, string(c))
		t.Log(n, n.dump(""))
	}
}

// func TestNodeLarge(t *testing.T) {
// 	data := make([]rune, 100)
// 	fill(data)
// 	n := newNode(data)

// 	n.insert(12, "krankelibrankelfnatt")
// 	t.Log(n, n.dump(""))
// 	n.erase(12, len("krankelibrankelfnatt"))
// 	t.Log(n, n.dump(""))
// }

func BenchmarkNodeSplit(b *testing.B) {
	b.StopTimer()
	data := make([]rune, 1024*256)
	fill(data)
	buf := newNodeEx(data, 4096)
	b.Log(buf.dump(""))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b2 := buf.clone()
		b2.split(512)
	}
}

var _ = rand.ExpFloat64

// func BenchmarkNodeInsertRand(b *testing.B) {
// 	r := rand.Perm(b.N)
// 	b.StopTimer()
// 	sdata := testinsert()
// 	buf := newNode(testbuffer().Runes())
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		buf.insert(r[i]%buf.length(), sdata)
// 	}
// }

func BenchmarkNodeInsertBegin(b *testing.B) {
	b.StopTimer()
	sdata := testinsert()
	buf := newNode(testbuffer().Runes())
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.insert(0, sdata)
	}
}

func BenchmarkNodeInsertMid(b *testing.B) {
	b.StopTimer()
	sdata := testinsert()
	buf := newNode(testbuffer().Runes())
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.insert(buf.length()/2, sdata)
	}
}

func BenchmarkNodeInsertEnd(b *testing.B) {
	b.StopTimer()
	sdata := testinsert()
	buf := newNode(testbuffer().Runes())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.insert(buf.length(), sdata)
	}
}
