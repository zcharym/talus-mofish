package autostart

import "testing"

func TestDefaultIdentifier(t *testing.T) {
	m := New("")
	if m.identifier != DefaultIdentifier {
		t.Fatalf("got %q, want %q", m.identifier, DefaultIdentifier)
	}
}
