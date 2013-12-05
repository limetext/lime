// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package render

const (
	Italic FontStyle = (1 << iota)
	Bold
	Underline
)

type (
	FontStyle int
	Font      struct {
		Name  string
		Size  float64
		Style FontStyle
	}

	FontMeasurement struct {
		Width, Height int
	}

	FontMetrics interface {
		Measure(Font, []rune) FontMeasurement
	}
)
