//go:build windows

package vdiupload

import (
	"context"
	"fmt"
)

// Unsupported backends are placeholders until Phase 2 Win32 work.
// They compile on Windows so the CLI can wire and dry-run the orchestrator.

type stubWindowLocator struct{}

func (stubWindowLocator) Find(ctx context.Context, match WindowMatch) (Window, error) {
	_ = ctx
	_ = match
	return Window{}, fmt.Errorf("%s: window locator not implemented (see DESIGN.md Phase 2)", CodeWindowNotFound)
}

type stubClipboard struct{}

func (stubClipboard) SetFiles(ctx context.Context, paths []string, cfg ClipboardConfig) (func(), error) {
	_ = ctx
	_ = paths
	_ = cfg
	return nil, fmt.Errorf("%s: clipboard not implemented (see DESIGN.md Phase 2)", CodeClipboardFailed)
}

type stubInput struct{}

func (stubInput) Focus(ctx context.Context, w Window) error {
	_ = ctx
	_ = w
	return fmt.Errorf("%s: input not implemented", CodeFocusFailed)
}

func (stubInput) RunKeys(ctx context.Context, keys []string) error {
	_ = ctx
	_ = keys
	return fmt.Errorf("%s: input not implemented", CodeInputFailed)
}

// DefaultDependencies returns Windows placeholder backends.
func DefaultDependencies() Dependencies {
	return Dependencies{
		Windows:   stubWindowLocator{},
		Clipboard: stubClipboard{},
		Input:     stubInput{},
	}
}
