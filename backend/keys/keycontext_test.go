package keys

import "testing"

func TestKeyContextUnmarshalError(t *testing.T) {
	var context KeyContext
	if err := context.UnmarshalJSON([]byte(``)); err == nil {
		t.Errorf("Expected error on loading empty string")
	}
}
