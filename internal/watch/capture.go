package watch

// WindowState describes whether the local VDI client can be captured.
type WindowState struct {
	Capturable bool
	Reason     string
	HWND       uintptr
	Width      int
	Height     int
}

// Capturer resolves VDI client windows and captures/OCRs screen regions.
type Capturer interface {
	Resolve(target Target) (WindowState, error)
	CaptureRegion(state WindowState, region Region) ([]byte, error)
	OCR(pixels []byte, width, height int) (string, error)
}
