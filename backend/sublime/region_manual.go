// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package sublime

import (
	"fmt"
	"github.com/quarnster/util/text"
	"lime/3rdparty/libs/gopy/lib"
)

func (o *Region) PyRichCompare(other py.Object, op py.Op) (py.Object, error) {
	if op != py.EQ && op != py.NE {
		return nil, fmt.Errorf("Can only do EQ and NE compares")
	}
	var o2 text.Region
	switch t := other.(type) {
	case *Region:
		o2 = t.data
	case *py.Tuple:
		if s := t.Size(); s != 2 {
			return nil, fmt.Errorf("Invalid tuple size: %d != 2", s)
		}
		if a, err := t.GetItem(0); err != nil {
			return nil, err
		} else if b, err := t.GetItem(1); err != nil {
			return nil, err
		} else if a2, ok := a.(*py.Long); !ok {
			return nil, fmt.Errorf("Can only compare with int tuples and other regions")
		} else if b2, ok := b.(*py.Long); !ok {
			return nil, fmt.Errorf("Can only compare with int tuples and other regions")
		} else {
			o2 = text.Region{int(a2.Int64()), int(b2.Int64())}
		}
	default:
		return nil, fmt.Errorf("Can only compare with int tuples and other regions")
	}
	if op == py.EQ {
		return toPython(o.data == o2)
	} else {
		return toPython(o.data != o2)
	}
}
