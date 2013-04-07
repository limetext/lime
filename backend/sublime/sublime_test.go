package sublime

import (
	//	"code.google.com/p/log4go"
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"lime/backend/primitives"
	"os"
	"path/filepath"
	"testing"
)

func TestSublime(t *testing.T) {
	ed := backend.GetEditor()
	ed.Console().Buffer().AddCallback(func(b *primitives.Buffer, pos, delta int) {
		fmt.Printf("%s", b.Data()[pos:pos+delta])
	})
	w := ed.NewWindow()
	w.NewView()
	Init()
	py.AddToPath("testdata")
	subl, err := py.Import("sublime")
	if err != nil {
		t.Fatal(err)
	}

	if w, err := _windowClass.Alloc(1); err != nil {
		t.Fatal(err)
	} else {
		(w.(*Window)).data = &backend.Window{}
		subl.AddObject("test_window", w)
	}
	if dir, err := os.Open("testdata"); err != nil {
		t.Error(err)
	} else if files, err := dir.Readdirnames(0); err != nil {
		t.Error(err)
	} else {
		for _, fn := range files {
			if filepath.Ext(fn) == ".py" {
				if _, err := py.Import(fn[:len(fn)-3]); err != nil {
					t.Error(err)
				} else {
					t.Log("Ran", fn)
				}
			}
		}
	}
}
