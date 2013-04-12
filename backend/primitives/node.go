package primitives

import (
	"fmt"
)

type (
	// Ropeish data structure.
	// http://en.wikipedia.org/wiki/Rope_(data_structure)
	// http://citeseer.ist.psu.edu/viewdoc/download?doi=10.1.1.14.9450&rep=rep1&type=pdf
	node struct {
		weight      int
		left, right *node
		data        []rune
	}
)

func (n *node) clone() *node {
	var lc, rc *node
	if n.left != nil {
		lc = n.left.clone()
	}
	if n.right != nil {
		rc = n.right.clone()
	}
	return &node{n.weight, lc, rc, n.data}
}

func (n *node) dump(indent string) string {
	indent += "\t"
	ret := fmt.Sprintf("%d, %s\n%sleft: ", n.weight, string(n.data), indent)
	if n.left == nil {
		ret += "nil\n"
	} else {
		ret += n.left.dump(indent)
	}
	ret += fmt.Sprintf("%sright: ", indent)
	if n.right == nil {
		ret += "nil\n"
	} else {
		ret += n.right.dump(indent)
	}
	return ret
}

func (n *node) Rune(pos int) rune {
	if s := pos - n.weight; s >= 0 && n.right != nil {
		return n.right.Rune(s)
	} else if n.weight > pos {
		if n.left != nil {
			return n.left.Rune(pos)
		} else {
			return n.data[pos]
		}
	}
	return ' '
}
func (n node) String() string {
	ret := ""
	if n.left != nil {
		ret += n.left.String()
	} else if len(n.data) != 0 {
		ret += string(n.data)
	}
	if n.right != nil {
		ret += n.right.String()
	}
	return ret
}

func (n *node) Substr(r Region) string {
	data := make([]rune, r.Size())
	b := r.Begin()
	for i := range data {
		data[i] = n.Rune(i + b)
	}
	return string(data)
}

func (n *node) find(pos int) *node {
	fmt.Println(pos)
	if n.weight < pos {
		return n.right.find(pos - n.weight)
	} else if n.left != nil {
		return n.left.find(pos)
	} else {
		return n
	}
}

func (n *node) simplify() {
	if (n.re()) && n.left != nil {
		*n = *n.left
	}
	if (n.le()) && n.right != nil {
		*n = *n.right
	}
	if n.empty() {
		n.weight = 0
		n.left = nil
		n.right = nil
	} else if n.weight < merge && (n.left != nil && n.left.leaf()) && (n.right != nil && n.right.leaf()) {
		n.left.join(n.right)
		*n = *n.left
	}
}

func (n *node) empty() bool {
	return len(n.data) == 0 && n.le() && (n.re())
}

func (n *node) re() bool {
	return n.right == nil || n.right.empty()
}

func (n *node) le() bool {
	return n.left == nil || n.left.empty()
}

func (n *node) leaf() bool {
	return n.le() && n.re()
}

var merge = 1024 * 4

func newNodeEx(data []rune, split int) *node {
	if len(data) > split {
		half := len(data) / 2
		return &node{half,
			newNodeEx(data[half:], split),
			newNodeEx(data[:half], split),
			nil,
		}
	}
	return &node{len(data), nil, nil, data}
}

func newNode(data []rune) *node {
	return newNodeEx(data, merge)
}

func (n *node) patch() {
	n.simplify()
	if n.left != nil {
		n.weight = n.left.length()
		if n.right != nil && n.right.left != nil && n.left.leaf() && n.right.left.leaf() && n.weight+n.right.weight < merge {
			r := n.right.split(n.right.weight)
			n.simplify()
			n.concat(r)
		}
	} else {
		n.weight = len(n.data)
	}
}

func (n *node) split(pos int) (right *node) {
	if n.weight < pos {
		return n.right.split(pos - n.weight)
	} else {
		if n.left != nil {
			right = n.left.split(pos)
		} else if n.right != nil {
			panic("shouldn't get here")
		} else {
			right = newNode(n.data[pos:])
			n.weight = pos
			n.data = n.data[:pos]
			return right
		}
	}
	if n.right != nil {
		right = &node{right.weight, right, n.right, nil}
	}
	n.right = nil
	n.patch()
	if right != nil {
		right.patch()
	}

	return right
}

func (n *node) length() int {
	ret := 0
	ret += n.weight
	if n.right != nil {
		ret += n.right.length()
	}
	return ret
}

func (n *node) join(other *node) {
	if len(n.data)+len(other.data) > merge {
		left := *n
		n.left = &left
		n.right = other
		n.data = nil
		n.weight = n.left.length()
	} else {
		// Allocating a new buffer as other nodes might have references
		// into sub positions in the original
		nd := make([]rune, 0, merge)
		n.data = append(nd, n.data...)
		n.data = append(n.data, other.data...)
		n.weight += other.weight
	}
}

func (n *node) concat(other *node) {
	if len(other.data) != 0 {
		if len(n.data) != 0 {
			//  If both arguments are short leaves, we produce a flat rope (leaf) consisting of the concatenation.
			n.join(other)
			return
		} else if n.right != nil {
			if n.right.le() {
				// If the left argument is a concatenation node whose right son is a short leaf,
				// and the right argument is also a short leaf,
				// then we concatenate the two leaves, and then concatenate the result to the left son of the left argument.
				n.right.concat(other)
				n.left.concat(n.right)
				n.right = nil
			} else {
				n.right.right.concat(other)
			}
		} else {
			n.right = other
		}
	} else {
		left := *n
		n.left = &left
		n.right = other
	}
	n.patch()
}

func (n *node) insert(position int, data string) {
	l := n.length()
	position = Clamp(0, l, position)
	r := []rune(data)
	left := newNode(r)
	if position >= l {
		n.concat(left)
	} else {
		right := n.split(position)
		n.concat(left)
		n.concat(right)
	}
}

func (n *node) erase(position, length int) {
	right := n.split(position + length)
	n.split(position)
	n.concat(right)
}
