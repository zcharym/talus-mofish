//go:build !windows

package vdiupload

import (
	"context"
	"fmt"
)

type unsupported struct{}

func (unsupported) Find(ctx context.Context, match WindowMatch) (Window, error) {
	return Window{}, fmt.Errorf("%s: vdiupload requires Windows", CodeUnsupportedOS)
}

func (unsupported) SetFiles(ctx context.Context, paths []string, cfg ClipboardConfig) (func(), error) {
	return nil, fmt.Errorf("%s: vdiupload requires Windows", CodeUnsupportedOS)
}

func (unsupported) Focus(ctx context.Context, w Window) error {
	return fmt.Errorf("%s: vdiupload requires Windows", CodeUnsupportedOS)
}

func (unsupported) RunKeys(ctx context.Context, keys []string) error {
	return fmt.Errorf("%s: vdiupload requires Windows", CodeUnsupportedOS)
}

// DefaultDependencies returns backends that fail clearly on non-Windows.
func DefaultDependencies() Dependencies {
	u := unsupported{}
	return Dependencies{Windows: u, Clipboard: u, Input: u}
}
