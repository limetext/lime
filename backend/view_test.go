// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"
	"github.com/limetext/lime/backend/textmate"
	"github.com/limetext/lime/backend/util"
	. "github.com/quarnster/util/text"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestView(t *testing.T) {
	var (
		w Window
		v = w.NewFile()
	)
	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	v.selection.Clear()
	r := []Region{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, r2 := range r {
		v.selection.Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range "1234" {
		for i := 0; i < v.selection.Len(); i++ {
			v.Insert(edit, v.selection.Get(i).Begin(), string(ins))
		}
	}
	v.EndEdit(edit)

	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "1234a1234b1234c1234d" {
		t.Error(d)
	}
	v.undoStack.Undo(true)
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "abcd" {
		t.Error("expected 'abcd', but got: ", d)
	}
	v.undoStack.Redo(true)
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", d)
	}

	v.selection.Clear()
	r = []Region{
		{0, 0},
		{5, 5},
		{10, 10},
		{15, 15},
	}
	for _, r2 := range r {
		v.selection.Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range []string{"hello ", "world"} {
		for i := 0; i < v.selection.Len(); i++ {
			v.Insert(edit, v.selection.Get(i).Begin(), ins)
		}
	}
	v.EndEdit(edit)

	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(d)
	}
	v.undoStack.Undo(true)

	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", d)
	}
	v.undoStack.Undo(true)
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "abcd" {
		t.Error("expected 'abcd', but got: ", d)
	}
	v.undoStack.Undo(true)
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "" {
		t.Error("expected '', but got: ", d)
	}
	v.undoStack.Redo(true)
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "abcd" {
		t.Error("expected 'abcd', but got: ", d)
	}

	v.undoStack.Redo(true)
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", d)
	}

	v.undoStack.Redo(true)
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(d)
	}
}

func TestErase(t *testing.T) {
	var (
		w Window
		v = w.NewFile()
		s = v.Sel()
	)
	edit := v.BeginEdit()
	v.Insert(edit, 0, "1234abcd5678abcd")
	v.EndEdit(edit)
	s.Clear()
	v.Sel().Add(Region{4, 8})
	v.Sel().Add(Region{12, 16})

	edit = v.BeginEdit()
	for i := 0; i < s.Len(); i++ {
		v.Erase(edit, s.Get(i))
	}
	v.EndEdit(edit)
	if !reflect.DeepEqual(s.Regions(), []Region{{4, 4}, {8, 8}}) {
		t.Error(s)
	}
	if d := v.buffer.Substr(Region{0, v.buffer.Size()}); d != "12345678" {
		t.Error(d)
	}
}

// This is not 100% what ST3 does
func TestExtractScope(t *testing.T) {
	var (
		w Window
		v = w.NewFile()
	)
	const (
		in      = "textmate/testdata/main.go"
		expfile = "testdata/scoperange.res"
		syntax  = "textmate/testdata/Go.tmLanguage"
	)
	v.Settings().Set("syntax", syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		t.Fatal(err)
	} else {
		//		v.rootNode = nil
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		last := Region{-1, -1}
		str := ""
		nr := Region{0, 0}
		for v.ExtractScope(1) == nr {
			time.Sleep(time.Millisecond)
		}
		for i := 0; i < v.buffer.Size(); i++ {
			if r := v.ExtractScope(i); r != last {
				str += fmt.Sprintf("%d (%d, %d)\n", i, r.A, r.B)
				last = r
			}
		}
		if d, err := ioutil.ReadFile(expfile); err != nil {
			if err := ioutil.WriteFile(expfile, []byte(str), 0644); err != nil {
				t.Error(err)
			}
		} else if diff := util.Diff(string(d), str); diff != "" {
			t.Error(diff)
		}
	}
}

// This is not 100% what ST3 does, but IMO ST3 is wrong
func TestScopeName(t *testing.T) {
	var (
		w Window
		v = w.NewFile()
	)
	const (
		in      = "textmate/testdata/main.go"
		expfile = "testdata/scopename.res"
		syntax  = "textmate/testdata/Go.tmLanguage"
	)
	v.Settings().Set("syntax", syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		t.Fatal(err)
	} else {
		//		v.rootNode = nil
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		last := ""
		str := ""
		lasti := 0
		for v.ScopeName(1) == "" {
			time.Sleep(250 * time.Millisecond)
		}
		for i := 0; i < v.buffer.Size(); i++ {
			if name := v.ScopeName(i); name != last {
				if last != "" {
					str += fmt.Sprintf("%d-%d: %s\n", lasti, i, last)
					lasti = i
				}
				last = name
			}
		}
		if i := v.Buffer().Size(); lasti != i {
			str += fmt.Sprintf("%d-%d: %s\n", lasti, i, last)
		}
		if d, err := ioutil.ReadFile(expfile); err != nil {
			if err := ioutil.WriteFile(expfile, []byte(str), 0644); err != nil {
				t.Error(err)
			}
		} else if diff := util.Diff(string(d), str); diff != "" {
			t.Error(diff)
		}

	}
}

func TestStress(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var (
		ed = GetEditor()
		w  = ed.NewWindow()
		v  = w.OpenFile("../frontend/termbox/main.go", 0)
	)
	syntax := "../3rdparty/bundles/go.tmbundle/Syntaxes/Go.tmLanguage"
	v.Settings().Set("syntax", syntax)
	for i := 0; i < 1000; i++ {
		e := v.BeginEdit()
		for i := 0; i < 100; i++ {
			v.Insert(e, 0, "h")
		}
		for i := 0; i < 100; i++ {
			v.Erase(e, Region{0, 1})
		}
		v.EndEdit(e)
	}
}

func TestTransform(t *testing.T) {
	sc, err := textmate.LoadTheme("../3rdparty/bundles/TextMate-Themes/GlitterBomb.tmTheme")
	if err != nil {
		t.Fatal(err)
	}
	var (
		w Window
		v = w.NewFile()
	)
	v.Settings().Set("syntax", "textmate/testdata/Go.tmLanguage")

	d, err := ioutil.ReadFile("view.go")
	if err != nil {
		t.Fatal(err)
	}
	e := v.BeginEdit()
	v.Insert(e, 0, string(d))
	v.EndEdit(e)

	time.Sleep(time.Second)
	a := v.Transform(sc, Region{0, 100}).Transcribe()
	v.Transform(sc, Region{100, 200}).Transcribe()
	c := v.Transform(sc, Region{0, 100}).Transcribe()
	if !reflect.DeepEqual(a, c) {
		t.Errorf("not equal:\n%v\n%v", a, c)
	}
}

func TestSaveAsNewFile(t *testing.T) {
	tests := []struct {
		text   string
		atomic bool
		file   string
	}{
		{
			"abc",
			false,
			"testdata/test",
		},
		{
			"abc",
			true,
			"testdata/test",
		},
	}
	for i, test := range tests {
		var (
			w Window
			v = w.NewFile()
			e = v.BeginEdit()
		)
		v.Settings().Set("atomic_save", test.atomic)
		v.Insert(e, 0, test.text)
		v.EndEdit(e)
		if err := v.SaveAs(test.file); err != nil {
			t.Fatalf("Test %d: Can't save to `%s`: %s", i, test.file, err)
		}
		data, err := ioutil.ReadFile(test.file)
		if err != nil {
			t.Fatalf("Test %d: Can't read `%s`: %s", i, test.file, err)
		}
		if string(data) != test.text {
			t.Errorf("Test %d: Expected `%s` contain %s, but got %s", i, test.file, test.text, data)
		}
		if err = os.Remove(test.file); err != nil {
			t.Errorf("Test %d: Couldn't remove test file %s", i, test.file)
		}
	}
}

func TestSaveAsOpenFile(t *testing.T) {
	var testfile string = "testdata/Default.sublime-settings"
	buf, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatalf("Can't read test file `%s`: %s", testfile, err)
	}
	tests := []struct {
		atomic bool
		as     string
	}{
		{
			true,
			"User.sublime-settings",
		},
		{
			true,
			"testdata/User.sublime-settings",
		},
		{
			true,
			"../User.sublime-settings",
		},
		{
			true,
			os.TempDir() + "/User.sublime-settings",
		},
		{
			false,
			"User.sublime-settings",
		},
		{
			false,
			"testdata/User.sublime-settings",
		},
	}
	for i, test := range tests {
		var (
			w Window
			v = w.OpenFile(testfile, 0)
		)
		v.Settings().Set("atomic_save", test.atomic)
		if err := v.SaveAs(test.as); err != nil {
			t.Fatalf("Test %d: Can't save to `%s`: %s", i, test.as, err)
		}
		if _, err := os.Stat(test.as); os.IsNotExist(err) {
			t.Fatalf("Test %d: The file `%s` wasn't created", i, test.as)
		}
		data, err := ioutil.ReadFile(test.as)
		if err != nil {
			t.Fatalf("Test %d: Can't read `%s`: %s", i, test.as, err)
		}
		if string(data) != string(buf) {
			t.Errorf("Test %d: Expected `%s` contain %s, but got %s", i, test.as, string(buf), data)
		}
		if err := os.Remove(test.as); err != nil {
			t.Errorf("Test %d: Couldn't remove test file %s", i, test.as)
		}
	}
}

func BenchmarkScopeNameLinear(b *testing.B) {
	var (
		w Window
		v = w.NewFile()
	)
	const (
		in     = "textmate/language_test.go"
		syntax = "textmate/testdata/Go.tmLanguage"
	)
	b.StopTimer()
	v.Settings().Set("syntax", syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		b.Fatal(err)
	} else {
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		b.StartTimer()
		for j := 0; j < b.N; j++ {
			for i := 0; i < v.buffer.Size(); i++ {
				v.ScopeName(i)
			}
		}
	}
}

func BenchmarkScopeNameRandom(b *testing.B) {
	var (
		w Window
		v = w.NewFile()
	)
	const (
		in     = "textmate/language_test.go"
		syntax = "textmate/testdata/Go.tmLanguage"
	)
	b.StopTimer()
	v.Settings().Set("syntax", syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		b.Fatal(err)
	} else {
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		p := rand.Perm(b.N)
		b.StartTimer()
		for _, i := range p {
			v.ScopeName(i)
		}
	}
}

func TestClassify(t *testing.T) {
	var w Window
	tests := []struct {
		text   string
		points []int
		expect []int
	}{
		{
			"",
			[]int{0, 10},
			[]int{3520, 3520},
		},
		{
			"abc Hi -test lime,te-xt\n\tclassify test-ing",
			[]int{0, 4, 5, 6, 7, 8, 13, 17, 18, 20, 21, 23, 24, 25, 34, 38, 39, 42},
			[]int{73, 49, 512, 2, 1028, 9, 1, 8198, 4105, 6, 9, 130, 64, 1, 1, 6, 9, 134},
		},
		{
			"(tes)ting cl][assify\n\npare(,,)nthe\\ses\n\t\n// Use",
			[]int{0, 4, 12, 13, 14, 20, 21, 22, 26, 27, 28, 29, 30, 34, 35, 39, 40, 41, 42, 43, 44, 47},
			[]int{5188, 8198, 8198, 12288, 4105, 130, 448, 65, 4102, 12288, 0, 12288, 8201, 6, 9, 64, 128, 1092, 0, 2056, 49, 134},
		},
	}
	for i, test := range tests {
		v := w.NewFile()
		e := v.BeginEdit()
		v.Insert(e, 0, test.text)
		v.EndEdit(e)
		for j, point := range test.points {
			if res := v.Classify(point); test.expect[j] != res {
				t.Errorf("Test %d: Expected %d from view.Classify(%d) but, got %d", i, test.expect[j], point, res)
			}
		}
	}
}
