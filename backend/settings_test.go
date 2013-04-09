package backend

import (
	"testing"
)

func TestSettings(t *testing.T) {
	var (
		s1, s2 HasSettings
		called bool
	)
	s1.Settings().SetParent(&s2)

	if v, ok := s1.Settings().Get("test", true).(bool); !ok || !v {
		t.Error(ok, v)
	}
	s2.Settings().Set("test", false)
	if v, ok := s1.Settings().Get("test", true).(bool); !ok || v {
		t.Error(ok, v)
	}

	s1.Settings().AddOnChange("something", func() {
		called = true
	})
	s2.Settings().Set("test", true)
	if !called {
		t.Error("Should have been called..")
	}
	called = false
	s1.Settings().ClearOnChange("something")
	s2.Settings().Set("test", true)
	if called {
		t.Error("Should not have been called..")
	}
}
