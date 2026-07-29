package watch

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPixels(pixels []byte) string {
	sum := sha256.Sum256(pixels)
	return hex.EncodeToString(sum[:])
}

type HashGate struct {
	last string
}

func (g *HashGate) Changed(pixels []byte) bool {
	h := HashPixels(pixels)
	if h == g.last {
		return false
	}
	g.last = h
	return true
}
