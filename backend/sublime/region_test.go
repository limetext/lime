package sublime

import (
	"github.com/qur/gopy/lib"
	"testing"
)

func TestRegion(t *testing.T) {
	py.AddToPath(".")
	if _, err := py.Import("region_test"); err != nil {
		t.Error(err)
	}
}
