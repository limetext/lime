package sublime

import (
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"os"
	"path/filepath"
	"testing"
)

func TestSublime(t *testing.T) {
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
