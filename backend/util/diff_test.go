// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		a, b, exp string
	}{
		{
			"test string",
			"test string",
			"",
		},

		{
			"test\nstring",
			"test\nstring",
			"",
		},

		{
			"test\r\nstring",
			"test\nstring",
			"",
		},

		{
			"test\n\rstring",
			"test\nstring",
			"  test\n- \rstring\n+ string",
		},
		{
			"test\nstring",
			"test\n\rstring",
			"  test\n- string\n+ \rstring",
		},

		{
			"t3st\nstring",
			"test\nstring",
			"- t3st\n+ test\n  string",
		},
		{
			"test\nstring",
			"t3st\nstring",
			"- test\n+ t3st\n  string",
		},

		{
			"t3st\nstr1ng",
			"test\nstring",
			"- t3st\n- str1ng\n+ test\n+ string",
		},
		{
			"test\nstring",
			"t3st\nstr1ng",
			"- test\n- string\n+ t3st\n+ str1ng",
		},

		{
			"t3st\nstr\n1ng",
			"test\nstr\ning",
			"- t3st\n+ test\n  str\n- 1ng\n+ ing",
		},
		{
			"test\nstr\ning",
			"t3st\nstr\n1ng",
			"- test\n+ t3st\n  str\n- ing\n+ 1ng",
		},

		{
			"t\ne\nst\nstr\n1ng",
			"t\ne\nst\nstr\ning",
			"  e\n  st\n  str\n- 1ng\n+ ing",
		},
		{
			"t\ne\nst\nstr\ning",
			"t\ne\nst\nstr\n1ng",
			"  e\n  st\n  str\n- ing\n+ 1ng",
		},

		{
			"t3\nst\nstr\nin\ng",
			"te\nst\nstr\nin\ng",
			"- t3\n+ te\n  st\n  str\n  in",
		},
		{
			"te\nst\nstr\nin\ng",
			"t3\nst\nstr\nin\ng",
			"- te\n+ t3\n  st\n  str\n  in",
		},

		{
			"t\n3st\nstring",
			"t\nest\nstring",
			"  t\n- 3st\n+ est\n  string",
		},
		{
			"t\nest\nstring",
			"t\n3st\nstring",
			"  t\n- est\n+ 3st\n  string",
		},

		{
			"t3\nst\nstr\n1ng",
			"te\nst\nstr\ning",
			"- t3\n+ te\n  st\n  str\n- 1ng\n+ ing",
		},
		{
			"te\nst\nstr\ning",
			"t3\nst\nstr\n1ng",
			"- te\n+ t3\n  st\n  str\n- ing\n+ 1ng",
		},

		{
			"T\ne\nst\ns\nt\nr\n1ng",
			"t\ne\nst\ns\nt\nr\ning",
			"- T\n+ t\n  e\n  st\n  s\n  t\n  r\n- 1ng\n+ ing",
		},
		{
			"t\ne\ns\nt\ns\nt\nr\ning",
			"T\ne\ns\nt\ns\nt\nr\n1ng",
			"- t\n+ T\n  e\n  s\n  t\n  s\n  t\n  r\n- ing\n+ 1ng",
		},

		{
			"T\ne\ns\nt\ns\nt\nr\n1ng",
			"t\ne\ns\nt\ns\nt\nr\ning",
			"- T\n+ t\n  e\n  s\n  t\n  s\n  t\n  r\n- 1ng\n+ ing",
		},
		{
			"t\ne\ns\nt\ns\nt\nr\ning",
			"T\ne\ns\nt\ns\nt\nr\n1ng",
			"- t\n+ T\n  e\n  s\n  t\n  s\n  t\n  r\n- ing\n+ 1ng",
		},

		{
			"T\ne\ns\nt\n\ns\nt\nr\n1ng",
			"t\ne\ns\nt\n\ns\nt\nr\ning",
			// TODO: No divider?
			"- T\n+ t\n  e\n  s\n  t\n  s\n  t\n  r\n- 1ng\n+ ing",
		},
		{
			"t\ne\ns\nt\n\ns\nt\nr\ning",
			"T\ne\ns\nt\n\ns\nt\nr\n1ng",
			// TODO: No divider?
			"- t\n+ T\n  e\n  s\n  t\n  s\n  t\n  r\n- ing\n+ 1ng",
		},

		{
			"teststr",
			"teststring",
			"- teststr\n+ teststring",
		},
		{
			"teststring",
			"teststr",
			"- teststring\n+ teststr",
		},

		{
			"test",
			"test\nstring",
			"  test\n+ string",
		},
		{
			"test\nstring",
			"test",
			"  test\n- string",
		},

		{
			"string",
			"test\nstring",
			"+ test\n  string",
		},
		{
			"test\nstring",
			"string",
			"- test\n  string",
		},

		{
			"str",
			"test\nstr\ning",
			"+ test\n  str\n+ ing",
		},
		{
			"test\nstr\ning",
			"str",
			"- test\n  str\n- ing",
		},

		{
			"ing",
			"test\nstr\ning",
			"+ test\n+ str\n  ing",
		},
		{
			"test\nstr\ning",
			"ing",
			"- test\n- str\n  ing",
		},

		{
			"test",
			"test\nstr\ning",
			"  test\n+ str\n+ ing",
		},
		{
			"test\nstr\ning",
			"test",
			"  test\n- str\n- ing",
		},
	}

	for i, test := range tests {
		out := Diff(test.a, test.b)

		if out != test.exp {
			t.Errorf("Test %d: Expected:\n%s\n\nBut got:\n%s", i, test.exp, out)
		}
	}
}
