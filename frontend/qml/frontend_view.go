// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"github.com/limetext/lime/backend"
	"github.com/limetext/lime/backend/render"
	"github.com/limetext/lime/backend/util"
	. "github.com/limetext/text"
	"gopkg.in/qml.v1"
	"io"
	"strings"
)

// A helper glue structure connecting the backend View with the qml code that
// then ends up rendering it.
type frontendView struct {
	bv            *backend.View
	qv            qml.Object
	FormattedLine []*lineStruct
	Title         lineStruct
}

// This allows us to trigger a qml.Changed on a specific line in the view so
// that only it is re-rendered by qml
type lineStruct struct {
	Text string
}

// htmlcol returns the hex color value for the given Colour object
func htmlcol(c render.Colour) string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

func (fv *frontendView) Line(index int) *lineStruct {
	return fv.FormattedLine[index]
}

func (fv *frontendView) Region(a int, b int) Region {
	return Region{a, b}
}

func (fv *frontendView) RegionLines() int {
	var count int = 0
	regs := fv.bv.Sel().Regions()
	if fv.bv.Buffer() != nil {
		for _, r := range regs {
			count += len(fv.bv.Buffer().Lines(r))
		}
	}
	return count
}

func (fv *frontendView) Setting(name string) interface{} {
	return fv.Back().Settings().Get(name, nil)
}

func (fv *frontendView) Back() *backend.View {
	return fv.bv
}

func (fv *frontendView) Fix(obj qml.Object) {
	fv.qv = obj

	for i := range fv.FormattedLine {
		_ = i
		obj.Call("addLine")
	}
}

func (fv *frontendView) bufferChanged(buf Buffer, pos, delta int) {
	prof := util.Prof.Enter("frontendView.bufferChanged")
	defer prof.Exit()

	row1, _ := buf.RowCol(pos)
	row2, _ := buf.RowCol(pos + delta)
	if row1 > row2 {
		row1, row2 = row2, row1
	}

	if delta > 0 && fv.qv != nil {
		r1 := row1
		if add := strings.Count(buf.Substr(Region{pos, pos + delta}), "\n"); add > 0 {
			nn := make([]*lineStruct, len(fv.FormattedLine)+add)
			copy(nn, fv.FormattedLine[:r1])
			copy(nn[r1+add:], fv.FormattedLine[r1:])
			for i := 0; i < add; i++ {
				nn[r1+i] = &lineStruct{Text: ""}
			}
			fv.FormattedLine = nn
			for i := 0; i < add; i++ {
				fv.qv.Call("insertLine", r1+i)
			}
		}
	}

	for i := row1; i <= row2; i++ {
		fv.formatLine(i)
	}
}

func (fv *frontendView) Erased(changed_buffer Buffer, region_removed Region, data_removed []rune) {
	fv.bufferChanged(changed_buffer, region_removed.B, region_removed.A-region_removed.B)
}

func (fv *frontendView) Inserted(changed_buffer Buffer, region_inserted Region, data_inserted []rune) {
	fv.bufferChanged(changed_buffer, region_inserted.A, region_inserted.B-region_inserted.A)
}

func (fv *frontendView) onChange(name string) {
	if name != "lime.syntax.updated" {
		return
	}
	// force redraw, as the syntax regions might have changed...
	for i := range fv.FormattedLine {
		fv.formatLine(i)
	}
}

func (fv *frontendView) formatLine(line int) {
	prof := util.Prof.Enter("frontendView.formatLine")
	defer prof.Exit()
	buf := bytes.NewBuffer(nil)
	vr := fv.bv.Buffer().Line(fv.bv.Buffer().TextPoint(line, 0))
	for line >= len(fv.FormattedLine) {
		fv.FormattedLine = append(fv.FormattedLine, &lineStruct{Text: ""})
		if fv.qv != nil {
			fv.qv.Call("addLine")
		}
	}
	if vr.Size() == 0 {
		if fv.FormattedLine[line].Text != "" {
			fv.FormattedLine[line].Text = ""
			t.qmlChanged(fv.FormattedLine[line], fv.FormattedLine[line])
		}
		return
	}
	recipie := fv.bv.Transform(scheme, vr).Transcribe()
	highlight_line := false
	if b, ok := fv.bv.Settings().Get("highlight_line", highlight_line).(bool); ok {
		highlight_line = b
	}
	lastEnd := vr.Begin()

	for _, reg := range recipie {
		if lastEnd != reg.Region.Begin() {
			fmt.Fprintf(buf, "<span>%s</span>", fv.bv.Buffer().Substr(Region{lastEnd, reg.Region.Begin()}))
		}
		fmt.Fprintf(buf, "<span style=\"white-space:pre; color:#%s; background:#%s\">%s</span>", htmlcol(reg.Flavour.Foreground), htmlcol(reg.Flavour.Background), fv.bv.Buffer().Substr(reg.Region))
		lastEnd = reg.Region.End()
	}
	if lastEnd != vr.End() {
		io.WriteString(buf, fv.bv.Buffer().Substr(Region{lastEnd, vr.End()}))
	}

	str := buf.String()

	if fv.FormattedLine[line].Text != str {
		fv.FormattedLine[line].Text = str
		t.qmlChanged(fv.FormattedLine[line], fv.FormattedLine[line])
	}
}
