package watch

import "testing"

func TestHashGate(t *testing.T) {
	gate := &HashGate{}
	pixels := []byte{1, 2, 3, 4}
	if !gate.Changed(pixels) {
		t.Fatal("first read should be changed")
	}
	if gate.Changed(pixels) {
		t.Fatal("same pixels should not be changed")
	}
	pixels[0] = 9
	if !gate.Changed(pixels) {
		t.Fatal("modified pixels should be changed")
	}
}
