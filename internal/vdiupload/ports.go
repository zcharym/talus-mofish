package vdiupload

import "context"

// Window is a resolved top-level VDI client window.
type Window struct {
	HWND   uintptr
	Title  string
	Width  int
	Height int
}

// WindowLocator finds the client window for a profile match rule.
type WindowLocator interface {
	Find(ctx context.Context, match WindowMatch) (Window, error)
}

// Clipboard stages host files for VDI clipboard redirect.
// Restore, if non-nil, should be deferred by the caller.
type Clipboard interface {
	SetFiles(ctx context.Context, paths []string, cfg ClipboardConfig) (restore func(), err error)
}

// InputAutomator focuses the client and runs profile steps that generate input.
type InputAutomator interface {
	Focus(ctx context.Context, w Window) error
	RunKeys(ctx context.Context, keys []string) error
}

// CompletionDetector waits until transfer success/failure UI is observed.
type CompletionDetector interface {
	WaitOCR(ctx context.Context, w Window, rule OCRWait) error
}
