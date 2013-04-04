package backend

import (
	"testing"
)

func TestSettings(t *testing.T) {
	var (
		s1, s2 HasSettings
	)
	s1.Settings().Parent = &s2

	if v, ok := s1.Settings().Get("test", true).(bool); !ok || !v {
		t.Error(ok, v)
	}
}
