package sublime

import (
	"github.com/qur/gopy/lib"
	"os"
	"path/filepath"
	"testing"
)

func TestSublime(t *testing.T) {
	py.AddToPath("testdata")
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
