package primitives

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestRowCol(t *testing.T) {
	var b Buffer
	if d, err := ioutil.ReadFile("./testdata/unittest.json"); err != nil {
		t.Fatal(err)
	} else {
		type Test struct {
			Offset          int
			Line            int
			Column          int
			LineUntilOffset string
			LineAtOffset    string
			WordAtOffset    string
		}
		var tests []Test
		if err := json.Unmarshal(d, &tests); err != nil {
			t.Fatal(err)
		} else {
			if d, err := ioutil.ReadFile("./testdata/unittest.cpp"); err != nil {
				t.Fatal(err)
			} else {
				b.Insert(0, string(d))
			}

			for i, test := range tests {
				var a Test
				a.Line, a.Column = b.RowCol(test.Offset)
				a.LineAtOffset = b.Substr(b.Line(test.Offset))
				a.WordAtOffset = b.Substr(b.Word(test.Offset))
				a.Offset = b.TextPoint(test.Line, test.Column)
				if a.Line != test.Line {
					t.Errorf("%d Line mismatch: %d != %d", i, a.Line, test.Line)
				}
				if a.Column != test.Column {
					t.Errorf("%d Column mismatch: %d != %d", i, a.Column, test.Column)
				}
				if a.Offset != test.Offset {
					t.Errorf("%d Offset mismatch: %d != %d", i, a.Offset, test.Offset)
				}
				if a.LineAtOffset != test.LineAtOffset {
					t.Errorf("%d LineAtOffset mismatch: '%s' != '%s'", i, a.LineAtOffset, test.LineAtOffset)
				}
				if a.WordAtOffset != test.WordAtOffset {
					t.Errorf("%d WordAtOffset mismatch: '%s' != '%s'", i, a.WordAtOffset, test.WordAtOffset)
				}
			}
		}
	}
	if r, c := b.RowCol(-1); r != 0 || c != 0 {
		t.Errorf("These should be 1 %d, %d", r, c)
	}
	if r, c := b.RowCol(b.Size() + 10); c != 0 {
		t.Errorf("Column should be 1 %d, %d", r, c)
	}
}
